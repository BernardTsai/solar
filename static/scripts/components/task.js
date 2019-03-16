Vue.component( 'task',
  {
    props:    ['model','view'],
    methods: {
      top: function(event) {
        if (event.Index1 < event.Index2) {
          return ((event.Index1)*view.automation.line + (event.Layer1+1)*view.automation.line/(event.Layers1+1)) + 'px'
        }
        return ((event.Index2)*view.automation.line + (event.Layer2+1)*view.automation.line/(event.Layers2+1)) + 'px'
      },
      height: function(event) {
        return Math.abs(
            ((event.Index1)*view.automation.line + (event.Layer1+1)*view.automation.line/(event.Layers1+1)) -
            ((event.Index2)*view.automation.line + (event.Layer2+1)*view.automation.line/(event.Layers2+1)) ) + 'px'
      },
      left: function(event) {
        return (100 + event.Time/this.model.Trace.Range*view.automation.width) + 'px'
      }
    },
    computed: {
      // range is the range beteen minimum and maximum
      range: function() {
        return this.model.Trace.Range
      },
      // events returns all trace related events
      events: function() {
        return this.model.Trace.Events
      },
      // tasks compiles all tasks from elements/clusters/instances
      tasks: function() {
        result = []

        // top level tasks
        Object.values(this.model.Trace.Tasks).forEach((task) => {
          result.push({
            UUID:      task.UUID,
            Started:   task.Started,
            Completed: task.Completed,
            Latest:    task.Latest,
            Status:    task.Status,
            State:     task.State,
            Version:   task.Version,
            Layer:     task.Layer,
            Layers:    task.Layers,
            Index:     0,
          })
        })

        // loop over elements
        Object.values(this.model.Trace.Elements).forEach((element) => {
          // add element tasks
          Object.values(element.Tasks).forEach((task) => {
            result.push({
              UUID:      task.UUID,
              Started:   task.Started,
              Completed: task.Completed,
              Latest:    task.Latest,
              Status:    task.Status,
              State:     task.State,
              Version:   task.Version,
              Layer:     task.Layer,
              Layers:    task.Layers,
              Index:     element.Index,
            })
          })

          // loop over clusters
          Object.values(element.Clusters).forEach((cluster) => {
            // add cluster tasks
            Object.values(cluster.Tasks).forEach((task) => {
              result.push({
                UUID:      task.UUID,
                Started:   task.Started,
                Completed: task.Completed,
                Latest:    task.Latest,
                Status:    task.Status,
                State:     task.State,
                Version:   task.Version,
                Layer:     task.Layer,
                Layers:    task.Layers,
                Index:     cluster.Index,
              })
            })

            // loop over instances
            Object.values(cluster.Instances).forEach((instance) => {
              // add instance tasks
              Object.values(instance.Tasks).forEach((task) => {
                result.push({
                  UUID:      task.UUID,
                  Started:   task.Started,
                  Completed: task.Completed,
                  Latest:    task.Latest,
                  Status:    task.Status,
                  State:     task.State,
                  Version:   task.Version,
                  Layer:     task.Layer,
                  Layers:    task.Layers,
                  Index:     instance.Index,
                })
              })
            })
          })
        })

        return result
      },
      // entities compiles all elements/clusters/instances
      entities: function() {
        result = []

        // loop over elements
        Object.values(this.model.Trace.Elements).forEach((element) => {
          result.push({
            ID:    element.Name,
            Name:  element.Name,
            Index: element.Index,
            Type:  "element"
          })

          // loop over clusters
          Object.values(element.Clusters).forEach((cluster) => {
            result.push({
              ID:    element.Name + "/" + cluster.Name,
              Name:  cluster.Name,
              Index: cluster.Index,
              Type:  "cluster"
            })

            // loop over instances
            Object.values(cluster.Instances).forEach((instance) => {
              result.push({
                ID:    element.Name + "/" + cluster.Name + "/" + instance.Name,
                Name:  instance.Name,
                Index: instance.Index,
                Type:  "instance"
              })
            })
          })
        })

        return result
      }
    },
    template: `
      <div id="task">
        <div class="entity"
          v-bind:style="{
            'top':         '0px',
            'height':      view.automation.line + 'px',
            'line-height': view.automation.line + 'px'}">
          <div class="label administrator"
            v-bind:style="{height: view.automation.line + 'px'}">
            Administrator
          </div>
        </div>
        <div class="entity"
          v-for="entity in entities"
          v-bind:style="{
            'top':         (entity.Index*view.automation.line) + 'px',
            'height':      view.automation.line + 'px',
            'line-height': view.automation.line + 'px'}">
          <div class="label"
            v-bind:class="entity.Type"
            v-bind:style="{height: view.automation.line + 'px'}">
            {{entity.Name}}
          </div>
        </div>
        <div class="task"
          v-bind:id="task.UUID"
          v-bind:title="task.UUID + ' - ' + task.Status"
          v-for="task in tasks"
          v-bind:style="{
            left:  (100 + task.Started/range*view.automation.width) + 'px',
            width: ((task.Latest-task.Started)/range*view.automation.width) + 'px',
            top:   (task.Index*view.automation.line + (task.Layer+1)*view.automation.line/(task.Layers+1)-2) + 'px'}">
        </div>
        <div class="task2"
          v-bind:id="task.UUID"
          v-bind:title="task.UUID + ' - ' + task.Status"
          v-for="task in tasks"
          v-bind:style="{
            left:  (100 + task.Started/range*view.automation.width) + 'px',
            width: ((task.Completed-task.Started)/range*view.automation.width) + 'px',
            top:   (task.Index*view.automation.line + (task.Layer+1)*view.automation.line/(task.Layers+1)-2) + 'px'}">
        </div>
        <div class="event"
          v-for="event in events"
          v-if="event.Task1 != event.Task2 && event.Task1 != ''"
          v-bind:class="event.Type"
          v-bind:id="event.UUID"
          v-bind:title="event.UUID + ' - ' + event.Type"
          v-bind:style="{
            top:    top(event),
            left:   left(event),
            height: height(event)}">
          <div class="arrow" v-bind:class="{down: event.Index1 < event.Index2, up: event.Index1 > event.Index2}"></div>
        </div>
        <div class="event2"
          v-for="event in events"
          v-bind:id="event.UUID"
          v-bind:class="event.Type"
          v-bind:title="event.UUID + ' - ' + event.Type"
          v-if="event.Task1 == event.Task2 || event.Task1 == ''"
          v-bind:style="{
            top:    top(event),
            left:   left(event)}">
        </div>
      </div>`
  }
)
