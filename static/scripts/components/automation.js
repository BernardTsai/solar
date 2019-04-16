Vue.component(
  'automation',
  {
    props: ['model', 'view'],
    methods: {
      // refresh reloads the current view
      refresh: function() {
        if (model.Trace) {
          loadTrace(this.view.domain, this.view.automation.task)
        } else {
          loadTasks(
            this.view.domain,
            this.view.automation.solution,
            this.view.automation.element,
            this.view.automation.cluster,
            this.view.automation.instance)
        }
      },
      // showTask displays the task trace
      showTask: function(task) {
        loadTrace(this.view.domain, task)
        this.view.automation.task = task
      },
      // hideTask displays the task selection window
      hideTask: function() {
        model.Trace               = null
        this.view.automation.task = null
      },
      selectSolution: function() {
        this.view.solution = this.view.automation.solution

        // load architectures
        if (this.view.domain != "" && this.view.solution != ""){
          loadSolution(this.view.domain, this.view.solution)
        }

        this.view.automation.element  = ""
        this.view.automation.cluster  = ""
        this.view.automation.instance = ""

        this.model.Trace = null

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

        this.model.Trace = null

        loadTasks(
          this.view.domain,
          this.view.automation.solution,
          this.view.automation.element,
          this.view.automation.cluster,
          this.view.automation.instance)
      },
      selectCluster: function() {
        view.automation.instance = ""

        this.model.Trace = null

        loadTasks(
          this.view.domain,
          this.view.automation.solution,
          this.view.automation.element,
          this.view.automation.cluster,
          this.view.automation.instance)
      },
      selectInstance: function() {
        this.model.Trace = null

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
          <div id="selector1" title="Solution">
            <strong>Context:</strong>
            <select v-model="view.automation.solution" v-on:change="selectSolution">
              <option selected value="">Solution</option>
              <option v-for="solution in model.Solutions">{{solution}}</option>
            </select>
          </div>

          <div id="selector2" v-if="model.Solution && view.automation.solution != ''" title="Element">
            <strong>/</strong>
            <select v-model="view.automation.element" v-on:change="selectElement">
              <option selected value="">Element</option>
              <option v-for="element in model.Solution.Elements">{{element.Element}}</option>
            </select>
          </div>

          <div id="selector3" v-if="model.Solution && view.automation.element != ''" title="Cluster">
            <strong>/</strong>
            <select v-model="view.automation.cluster" v-on:change="selectCluster">
              <option selected value="">Cluster</option>
              <option v-for="cluster in model.Solution.Elements[view.automation.element].Clusters">{{cluster.Version}}</option>
            </select>
          </div>

          <div  id="selector4" v-if="model.Solution && view.automation.element != '' && view.automation.cluster != ''" title="Instance">
            <strong>/</strong>
            <select v-model="view.automation.instance" v-on:change="selectInstance">
              <option selected value="">Instance</option>
              <option v-for="instance in model.Solution.Elements[view.automation.element].Clusters[view.automation.cluster].Instances">{{instance.UUID}}</option>
            </select>
          </div>

          <div  id="selector5" v-if="model.Trace" title="Task">
            <strong>/</strong>
            <select>
              <option selected>{{this.view.automation.task}}</option>
            </select>
          </div>

          <div id="hide"    @click="hideTask" v-if="model.Trace"       title="hide"><i class="fas fa-eye-slash"></i></div>
          <div id="refresh" @click="refresh"  v-if="view.solution!=''" title="refresh"><i class="fas fa-recycle"></i></div>

        </div>

        <table id="tasks" v-if="!model.Trace">
          <tr>
            <th>Type</th>
            <th>UUID</th>
            <th>Version/State</th>
            <th>Started</th>
            <th>Completed</th>
            <th>Latest</th>
            <th>Status</th>
          </tr>
          <tr v-for="task in model.Tasks" @click="showTask(task.UUID)">
            <td>{{task.Type}}</td>
            <td>{{task.UUID}}</td>
            <td v-if="task.Type != 'InstanceTask'">{{task.Version}}</td>
            <td v-if="task.Type == 'InstanceTask'">{{task.State}}</td>
            <td>{{task.Started}}</td>
            <td>{{task.Completed}}</td>
            <td>{{task.Latest}}</td>
            <td class="status" :class="task.Status">{{task.Status}}</td>
          </tr>
          <tr v-if="model.Tasks.length == 0">
            <td colspan=7 class="noentries">no tasks</td>
          </tr>
        </table>
        <task v-if="model.Trace"
          v-bind:model="model"
          v-bind:view="view"></task>
      </div>`
  }
)
