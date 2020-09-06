// loadData synchronously retrieves data from an URL as string
var loadData = async(url) => {
  response = await fetch(url);
  text     = await response.text();

  return text
}

//------------------------------------------------------------------------------

// uuid generates a new universal unique ID
function uuid() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = Math.random()*16|0, v = c === 'x' ? r : (r&0x3|0x8);
        return v.toString(16);
    });
}

//------------------------------------------------------------------------------

// getName determines name part of a label
function getName(label) {
  pos     = label.lastIndexOf(" - V")
  name    = label.substring(0, pos)

  return name
}

//------------------------------------------------------------------------------

// getVersion determines version part of a label
function getVersion(label) {
  pos     = label.lastIndexOf(" - V")
  version = label.substring(pos+3)

  return version
}

//------------------------------------------------------------------------------

// dump prints an object as yaml to the console
function dump(obj) {
  console.log(jsyaml.safeDump(obj))
}

//------------------------------------------------------------------------------

// validateSchema validates a schema definition
function validateSchema(schemaString) {
  schemaObject = null

  // empty schemas are always valid
  if (schemaString == "") {
    return ""
  }

  // convert schema string to schema object
  try {
    schemaObject = jsyaml.safeLoad(schemaString)
  } catch (err) {
    return err.message
  }

  // create validator from schema object
  try {
    validator = new Ajv().compile(schemaObject)
  } catch(err) {
    return err.message
  }

  return ""
}

//------------------------------------------------------------------------------

// validateParameters validates parameters against a schema definition
function validateParameters(schemaString, parametersString) {
  schemaObject = null
  paramsObject = null

  // empty schemas and parameters are always valid
  if (schemaString == "" && parametersString == "") {
    return ""
  }

  // convert yaml parameters string to parameters object
  try {
    paramsObject = jsyaml.safeLoad(parametersString)
  } catch (err) {
    return err.message
  }

  // empty schemas are always valid
  if (schemaString == "") {
    return ""
  }

  // convert schema string to schema object
  try {
    schemaObject = jsyaml.safeLoad(schemaString)
  } catch (err) {
    return err.message
  }

  // create validator from schema object
  try {
    validator = new Ajv().compile(schemaObject)
  } catch(err) {
    return err.message
  }

  // validate parameters object
  valid = validator(paramsObject)
  if (!valid) {
    return "Parameters are invalid"
  }

  return ""
}

//------------------------------------------------------------------------------
