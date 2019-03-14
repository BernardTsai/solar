Vue.component(
  'automation',
  {
    props: ['model', 'view'],
    methods: {
      selectSolution: function() {
        this.view.solution = this.view.automation.solution

        // load architectures
        if (this.view.domain != "" && this.view.solution != ""){
          loadSolution(this.view.domain, this.view.solution)
        }

        this.view.automation.element  = ""
        this.view.automation.cluster  = ""
        this.view.automation.instance = ""

        loadTasks(
          this.view.domain,
          this.view.automation.solution,
          this.view.automation.element,
          this.view.automation.cluster,
          this.view.automation.instance)
      },
      selectElement: function() {
        view.automation.cluster  = ""
        view.automation.instance = ""

        loadTasks(
          this.view.domain,
          this.view.automation.solution,
          this.view.automation.element,
          this.view.automation.cluster,
          this.view.automation.instance)
      },
      selectCluster: function() {
        view.automation.instance = ""

        loadTasks(
          this.view.domain,
          this.view.automation.solution,
          this.view.automation.element,
          this.view.automation.cluster,
          this.view.automation.instance)
      },
      selectInstance: function() {
        loadTasks(
          this.view.domain,
          this.view.automation.solution,
          this.view.automation.element,
          this.view.automation.cluster,
          this.view.automation.instance)
      }
    },
    template: `
      <div id="automation">
        <div id="selector">
          <div id="selector1">
            <strong>Solution:</strong>
            <select v-model="view.automation.solution" v-on:change="selectSolution">
              <option disabled selected value="">Please select one</option>
              <option v-for="solution in model.Solutions">{{solution}}</option>
            </select>
          </div>

          <div id="selector2" v-if="model.Solution && view.automation.solution != ''">
            <strong>Element:</strong>
            <select v-model="view.automation.element" v-on:change="selectElement">
              <option disabled selected value="">Please select one</option>
              <option v-for="element in model.Solution.Elements">{{element.Element}}</option>
            </select>
          </div>

          <div id="selector3" v-if="model.Solution && view.automation.element != ''">
            <strong>Cluster:</strong>
            <select v-model="view.automation.cluster" v-on:change="selectCluster">
              <option disabled selected value="">Please select one</option>
              <option v-for="cluster in model.Solution.Elements[view.automation.element].Clusters">{{cluster.Version}}</option>
            </select>
          </div>

          <div  id="selector4" v-if="model.Solution && view.automation.element != '' && view.automation.cluster != ''">
            <strong>Instance:</strong>
            <select v-model="view.automation.instance" v-on:change="selectInstance">
              <option disabled selected value="">Please select one</option>
              <option v-for="instance in model.Solution.Elements[view.automation.element].Clusters[view.automation.cluster].Instances">{{instance.UUID}}</option>
            </select>
          </div>

        </div>

        <table id="tasks">
          <tr>
            <th>Type</th>
            <th>UUID</th>
            <th>State</th>
            <th>Started</th>
            <th>Completed</th>
            <th>Latest</th>
            <th>Status</th>
          </tr>
          <tr v-for="task in model.Tasks">
            <td>{{task.Type}}</td>
            <td>{{task.UUID}}</td>
            <td>{{task.State}}</td>
            <td>{{task.Started}}</td>
            <td>{{task.Completed}}</td>
            <td>{{task.Latest}}</td>
            <td class="status" :class="task.Status">{{task.Status}}</td>
          </tr>
        </table>
      </div>`
  }
)
