package util

import (
  "testing"
  "os"
)

//------------------------------------------------------------------------------

// TestConvertFromYAML tests the ConvertFromYAML function.
func TestConvertFromYAML(t *testing.T) {
  var object map[string]int

  yaml := "a: 2"
  err := ConvertFromYAML(yaml, &object)
  if err != nil {
    t.Errorf("ConvertFromYAML should have been able to convert yaml:\n%s", err)
  }

  yaml = "a: 2:a"
  err = ConvertFromYAML(yaml, &object)
  if err == nil {
    t.Errorf("ConvertFromYAML should have complained about invalid yaml input")
  }
}

//------------------------------------------------------------------------------

// TestConvertToYAML tests the ConvertToYAML function.
func TestConvertToYAML(t *testing.T) {
  object := map[string]int{"a":2}

  _, err := ConvertToYAML(&object)
  if err != nil {
    t.Errorf("ConvertToYAML should have been able to convert yaml:\n%s", err)
  }
}

//------------------------------------------------------------------------------

// TestConvertYAMLToJSON tests the ConvertYAMLToJSON function.
func TestConvertYAMLToJSON(t *testing.T) {
  input1 := `
    a: 1
    b: 'test'
    c:
      - 2
      - 3
      - 4
    3: 'test2'`

  _, err := ConvertYAMLToJSON( []byte(input1) )
  if err != nil {
    t.Errorf("ConvertYAMLToJSON should have been able to convert input 1")
  }
}

//------------------------------------------------------------------------------

// TestSaveAndLoadFile tests the SaveFile and LoadFile functions.
func TestLoadFile(t *testing.T) {
  filename1 := "./util_test.yaml"
  content1  := "test"

  // cleanup routine
  defer func() {
    os.Remove(filename1)
  }()

  // create configuration file
  err:= SaveFile(filename1, content1)
  if err != nil {
    t.Fatalf("Unable to save yaml file:\n%s", err )
  }

  _, err = LoadFile(filename1)
  if err != nil {
    t.Fatalf("Unable to load yaml file:\n%s", err )
  }
}

//------------------------------------------------------------------------------

// TestLoadYAML tests the LoadYAML function.
func TestLoadYAML(t *testing.T) {
  var object map[string]int

  filename1 := "./util_test.yaml"
  yaml1     := "a: 2"
  filename2 := "./util_test.yml"
  yaml2     := "  a 2\nb"
  filename3 := "./util_test.yiml"

  // cleanup routine
  defer func() {
    os.Remove(filename1)
    os.Remove(filename2)
  }()

  // create configuration file
  f, err := os.Create(filename1)
  if err != nil {
    t.Fatalf("Unable to create yaml file:\n%s", err )
  }

  // write content
  _, err = f.WriteString(yaml1)
  if err != nil {
    t.Fatalf("Unable to write yaml file:\n%s", err )
  }

  // close configuration file
  err = f.Close()
  if err != nil {
    t.Fatalf("Unable to write yaml file:\n%s", err )
  }

  // read yaml from existing file
  err = LoadYAML(filename1, &object)
  if err != nil {
    t.Errorf("Unable to read and parse yaml file:\n%s", err)
  }

  // create configuration file
  f, err = os.Create(filename2)
  if err != nil {
    t.Fatalf("Unable to create yaml file:\n%s", err )
  }

  // write content
  _, err = f.WriteString(yaml2)
  if err != nil {
    t.Fatalf("Unable to write yaml file:\n%s", err )
  }

  // close configuration file
  err = f.Close()
  if err != nil {
    t.Fatalf("Unable to write yaml file:\n%s", err )
  }

  // read yaml from existing file
  err = LoadYAML(filename2, &object)
  if err == nil {
    t.Errorf("LoadYAML should have reported an error regarding invalid data structure")
    DumpYAML(object)
  }

  // read yaml from non existing file
  err = LoadYAML(filename3, &object)
  if err == nil {
    t.Errorf("LoadYAML should have reported an error regarding missing file")
  }
}

//------------------------------------------------------------------------------

// TestSaveYAML tests the SaveYAML function.
func TestSaveYAML(t *testing.T) {
  filename1 := "./util_test.yaml"
  object1   := map[string]int{"a":2}
  filename2 := "./unknown_directory/util_test.yaml"

  // cleanup routine
  defer func() {
    os.Remove(filename1)
  }()

  // save yaml from non existing file
  err := SaveYAML(filename1, &object1)
  if err != nil {
    t.Errorf("SaveYAML failed:\n%s", err)
  }

  // try to save yaml twice
  err = SaveYAML(filename2, &object1)
  if err == nil {
    t.Errorf("SaveYAML should have failed when attempting to save entity")
  }
}

//------------------------------------------------------------------------------

// TestMisc tests the other miscellaneous functions of the util package.
func TestMisc(t *testing.T) {
  UUID()
  DumpYAML("")
  Print("")
}

//------------------------------------------------------------------------------
