Vue.component(
  'administration',
  {
    props: ['model', 'view'],
    methods: {
      // clearModel deletes the current model
      clearModel: function() {
        resetModel()
        .then(() => loadDomains())
        .then(() => this.view.domain = "")
        .then(() => alert("Model has been reset"))
      },
      // importModel imports a new model from an uploaded file
      importModel: function() {
        uploadButton = document.getElementById("file")
        uploadButton.click()
      },
      // uploadModel uploads a model to the repository
      uploadModel: function(event) {
        file = event.target.files[0]

        reader = new FileReader();

        reader.onload = ((e) => {
          saveModel(e.target.result)
          .then(() => loadDomains())
          .then(() => alert("Model has been imported"))
        })

        // Klartext mit Zeichenkodierung UTF-8 auslesen.
        reader.readAsText(file, "utf8")
      },
      // exportModel downloads the existing model
      exportModel: function() {
        window.open("/model", "_blank")
      },
      // createDomain creates a new domain
      createDomain: function() {
        // ask for name of new domain and add a new domain
        name = prompt("Name of the new domain:")
        if (name != null && name != "" && name != "null") {
          saveDomain(name)
          .then(() => loadDomains())
          .then(() => {
            view.domain = name
            navComponents()
            selectDomain()
          })
        }
      },
      // removeDomain removes an existing domain
      removeDomain: function() {
        // ask for name of an existing domain and remove it
        name = prompt("Name of the domain to be removed:")
        if (name != null && name != "" && name != "null") {
          deleteDomain(name)
          .then(() => loadDomains())
          .then(() => alert("Domain has been removed"))
          .then(() => {
            view.domain = ""
          })
        }
      }
    },
    template: `
      <div id="administration" v-if="view.nav=='Administration'">

        <div id="selector">
        </div>

        <div id="actions">
          <div class="action" @click="createDomain">
            <i class="fas fa-plus-square fa-lg"></i>&nbsp;Create Domain
          </div>

          <div class="action" @click="removeDomain">
            <i class="fas fa-minus-square fa-lg"></i>&nbsp;Remove Domain
          </div>

          <div class="action" @click="clearModel">
            <i class="fas fa-window-close fa-lg"></i>&nbsp;Reset Model
          </div>

          <div class="action" @click="importModel">
            <i class="fas fa-caret-square-up fa-lg"></i>&nbsp;Import Model
            <input type="file" id="file" name="model" @change="uploadModel"/>
          </div>

          <div class="action" @click="exportModel">
            <i class="fas fa-caret-square-down fa-lg"></i>&nbsp;Export Model
          </div>
        </div>

      </div>`
  }
)
