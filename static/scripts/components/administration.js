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
        .then(() => this.view.modelDomain = "")
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
      // addDomain creates a new domain
      addDomain: function() {
        // ask for name of new domain and add a new domain
        name = prompt("Name of the new domain:")
        if (name != null && name != "" && name != "null") {
          saveDomain(name)
          .then(() => loadDomains())
        }
      },
      // selectDomain selects a domain for editing
      selectDomain: function(domain) {
        view.domain = domain

        // load catalog
        if (view.domain == ""){
          model.Catalog = []
        } else {
          loadCatalog(view.domain)
        }

        // load components
        if (view.domain == ""){
          model.Components = []
        } else {
          loadComponents(view.domain)
        }

        // load architectures
        if (view.domain == ""){
          model.Architectures = []
          model.Architecture  = null
        } else {
          loadArchitectures(view.domain)
          model.Architecture  = null
        }
        view.architecture = ""

        // load solutions
        if (view.domain == ""){
          model.Solutions = []
          model.Solution  = null
        } else {
          loadSolutions(view.domain)
          model.Solution  = null
        }
        view.solution = ""

        // activate administration if no domain has been selected
        if (view.domain == "") {
          navAdministration()
        } else if (view.nav == "" || view.nav == "Administration") {
          view.nav = "Components"
        }
      },
      // selectModelDomain selects a specific model domain from view.modelDomain
      selectModelDomain: function(domain) {
        this.view.modelDomain = domain
        loadControllers(view.modelDomain)
      },
      // removeDomain removes an existing domain
      removeDomain: function(domain) {
        // ask for name of an existing domain and remove it
        deleteDomain(domain)
        .then(() => loadDomains())
        .then(() => {
          view.modelDomain = ""
        })
      },
      // addController adds a new controller
      addController: function() {
        // ask for name of the new controller
        name = prompt("Image name of the controller to be added:")
        if (name != null && name != "" && name != "null") {
          controller = {
            Controller:  uuid(),
            Version:     "V0.0.0",
            Image:       name,
            URL:         "",
            Status:      ""
          }

          addController(view.modelDomain, controller)
          .then(() => loadControllers(view.modelDomain))
        }
      },
      // deleteController removes an existing controller
      deleteController: function(controller) {
        if (controller.Image != "") {
          deleteController(view.modelDomain, controller)
          .then(() => loadControllers(view.modelDomain))
        }
      }
    },
    template: `
      <div id="Administration" v-if="view.nav=='Administration'">

        <!-- model domain selector -->
        <div id="selector">
          <strong>Model</strong>
          <div class="buttons">
            <button class="action" title="import model" @click="importModel">
              Import&nbsp;<i class="fas fa-arrow-circle-up"></i>
              <input type="file" id="file" name="model" @change="uploadModel" style=""/>
            </button>
            <button class="action" title="export model" @click="exportModel">
              Export&nbsp;<i class="fas fa-arrow-circle-down"></i>
            </button>
            <button class="action" title="reset model" @click="clearModel">
              Reset&nbsp;<i class="fas fa-times-circle"></i>
            </button>
          </div>
        </div>

        <!-- model domains -->
        <div id="modelDomains">
          <div class="header">
            <h3>Domains:</h3>
          </div>

          <table class="components">
            <thead>
              <tr>
                <th>Domain</th>
                <th @click="addDomain()" title="add new domain">
                  <i class="fas fa-plus-circle"></i>
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="domain in model.Domains">
                <td @click="selectModelDomain(domain)" title="view domain">{{domain}}</td>
                <td  @click="selectDomain(domain)" title="select domain">
                  <i class="fas fa-edit"></i>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- domain details -->
        <div id="domainDetails" v-if="view.modelDomain!=''">
          <div class="header">
            <h3>Domain: {{view.modelDomain}}</h3>
            <button class="modal-default-button" v-on:click="removeDomain(view.modelDomain)" title="delete domain">
              Delete <i class="fas fa-times-circle">
            </button>
            <button class="modal-default-button" v-on:click="selectModelDomain(view.modelDomain)" title="refresh">
              Refresh <i class="fas fa-recycle">
            </button>
            <button class="modal-default-button" v-on:click="selectDomain(view.modelDomain)" title="select domain">
              Edit <i class="fas fa-edit">
            </button>
          </div>

          <table style="width: 100%">
            <col width="10*">
            <col width="990*">
            <tr>
              <td><strong>Controllers:</strong></td>
              <td>

                <table id="controllers">
                  <thead>
                    <tr>
                      <th>Controller</th>
                      <th>Version</th>
                      <th>Image</th>
                      <th>URL</th>
                      <th>Status</th>
                      <th class="center" @click="addController" title="add controller">
                        <i class="fas fa-plus-circle"></i>
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="controller in model.Controllers">
                      <td>{{controller.Controller}}</td>
                      <td>{{controller.Version}}</td>
                      <td>{{controller.Image}}</td>
                      <td>{{controller.URL}}</td>
                      <td>{{controller.Status}}</td>
                      <td @click="deleteController(controller)" title="delete controller">
                        <div v-if="controller.Image!=''"><i class="fas fa-minus-circle"></i></div>
                      </td>
                    </tr>
                  </tbody>
                </table>

              </td>
            </tr>
          </table>

        </div>
      </div>`
  }
)
