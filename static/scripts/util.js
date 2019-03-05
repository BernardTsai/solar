var loadData = async(url) => {
  response = await fetch(url);
  text     = await response.text();

  return text
}

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
