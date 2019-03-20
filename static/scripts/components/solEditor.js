Vue.component(
  'solEditor',
  {
    props: ['model', 'view','element'],
    methods: {
      // showTasks switches to the corresponding automation tasks
      showTasks: function() {
        model.Tasks = []
        model.Trace = null

        view.automation.solution = view.solution
        view.automation.element  = this.element.Element
        view.automation.cluster  = ""
        view.automation.instance = ""

        loadTasks(
          view.domain,
          view.automation.solution,
          view.automation.element,
          view.automation.cluster,
          view.automation.instance)

        view.nav = "Automation"
      },
      // editConfiguration opens editor for editing a configuration
      editConfiguration: function(rel) {
        if (rel != "") {
          this.view.se.ConfTitle     = "Configuration of relationship: " + rel.Relationship
          this.view.se.Configuration = rel.Configuration
          this.view.se.Dependency    = rel.Dependency
        } else {
          this.view.se.ConfTitle     = "Configuration of element: "
          this.view.se.Configuration = this.model.Element.Configuration
          this.view.se.Dependency    = ""
        }
        this.$forceUpdate()
      },
      // discardConfiguration closes the configuration editor
      discardConfiguration: function() {
        this.view.se.Configuration = null
        this.view.se.Dependency    = ""

        this.$forceUpdate()
      }
    },
    template: `
      <div class="solEditor">
        <div class="header">
          <h3>Element: {{element.Element}} ({{element.Component}})</h3>
          <button class="modal-default-button" @click="showTasks()">
            Tasks <i class="fas fa-cogs">
          </button>
        </div>

        <table style="width: 100%">
          <col width="10*">
          <col width="990*">
          <tr>
            <td><strong>Element:</strong></td>
            <td>
              <input type="text" readonly
                v-model="element.Element">
              </input>
            </td>
          </tr>
          <tr>
            <td>&nbsp;Component:</td>
            <td>
              <input type="text" readonly
                v-model="element.Component">
              </input>
            </td>
          </tr>
          <tr>
            <td>&nbsp;Configuration:</td>
            <td>
              <textarea id="configuration" rows=5 readonly
                @click="editConfiguration('')"
                v-model="element.Configuration">
              </textarea>
            </td>
          </tr>

          <tr>
            <td><strong>Clusters:</strong></td>
            <td>
              <select
                v-model="view.se.Cluster">
                <option value="">please select</option>
                <option v-for="cluster in element.Clusters">{{cluster.Version}}</option>
              </select>
            </td>
          </tr>
          <tr v-if="view.se.Cluster != ''">
            <td>&nbsp;Target State:</td>
            <td>
              <input type="text" readonly v-model="element.Clusters[view.se.Cluster].Target"/>
            </td>
          </tr>
          <tr v-if="view.se.Cluster != ''">
            <td>&nbsp;Current State:</td>
            <td>
              <input type="text" readonly v-model="element.Clusters[view.se.Cluster].State"/>
            </td>
          </tr>
          <tr v-if="view.se.Cluster != ''">
            <td>&nbsp;Minimum Size:</td>
            <td>
              <input type="text" readonly v-model="element.Clusters[view.se.Cluster].Min"/>
            </td>
          </tr>
          <tr v-if="view.se.Cluster != ''">
            <td>&nbsp;Maximum Size:</td>
            <td>
              <input type="text" readonly v-model="element.Clusters[view.se.Cluster].Max"/>
            </td>
          </tr>
          <tr v-if="view.se.Cluster != ''">
            <td>&nbsp;Current Size:</td>
            <td>
              <input type="text" readonly v-model="element.Clusters[view.se.Cluster].Size"/>
            </td>
          </tr>
          <tr v-if="view.se.Cluster != ''">
            <td>&nbsp;Configuration:</td>
            <td>
              <textarea id="configuration2" rows=5 readonly
                @click="editConfiguration('')"
                v-model="element.Clusters[view.se.Cluster].Configuration">
              </textarea>
            </td>
          </tr>

          <tr v-if="view.se.Cluster != ''">
            <td><strong>Relationships:</strong></td>
            <td>

              <table class="relationships">
                <thead>
                  <tr>
                    <th>Relationship</th>
                    <th>Dependency</th>
                    <th>Element</th>
                    <th>Configuration</th>
                </thead>
                <tbody>
                  <tr v-for="rel in element.Clusters[view.se.Cluster].Relationships">
                    <td>
                      <input type="text" readonly v-model="rel.Relationship"/>
                    </td>
                    <td>
                      <input type="text" readonly v-model="rel.Dependency"/>
                    </td>
                    <td>
                      <input type="text" readonly v-model="rel.Element"/>
                    </td>
                    <td @click="editConfiguration(rel)">
                      <i class="fas fa-edit"></i>
                    </td>
                  </tr>
                </tbody>
              </table>

            </td>
          </tr>
        </table>

        <div class="configurationEditor" v-if="view.se.Configuration != null">
          <div class="modal">
            <h3>{{view.se.ConfTitle}}</h3>
            <textarea v-focus @keyup.esc="discardConfiguration()" v-model="view.se.Configuration"></textarea>
            <button class="modal-default-button" @click="discardConfiguration()">
              Cancel <i class="fas fa-times-circle">
            </button>

          </div>
        </div>

      </div>`
  }
)
