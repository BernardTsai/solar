Vue.component(
  'node_destination',
  {
    props: ['model', 'view', 'node', 'destination', 'idx'],
    computed: {
      cx: function() {
        var portIndex = parseInt(this.idx) + 1
        var ports     = this.node.Destinations.length + 1
        var radius    = this.view.graph.port.diameter/2
        var w         = this.view.graph.node.width - 2 * radius

        return (radius + portIndex/ports*w)
      },
      tag: function() {
        return this.destination.Tag
      }
    },
    template: `
      <circle class="destination"
        :class="destination.State"
        :cx="cx"
        cy=0
        r=4>
        <title>Cluster: {{tag}}</title>
      </circle>`
  }
)

//------------------------------------------------------------------------------

Vue.component(
  'node_source',
  {
    props: ['model', 'view', 'node', 'source', 'idx'],
    computed: {
      cx: function() {
        var portIndex = parseInt(this.idx) + 1
        var ports     = this.node.Sources.length + 1
        var radius    = this.view.graph.port.diameter/2
        var w         = this.view.graph.node.width - 2 * radius

        return (radius + portIndex/ports*w)
      },
      tag: function() {
        return this.source.Tag
      }
    },
    template: `
      <circle class="source"
        :cx="cx"
        cy=40
        r=4>
        <title>Dependency: {{tag}}</title>
      </circle>`
  }
)

//------------------------------------------------------------------------------

Vue.component(
  'node',
  {
    props: ['model', 'view', 'node'],
    computed: {
      x: function() {
        var w   = this.view.graph.node.width
        var dx  = this.view.graph.dx
        var col = this.node.Column

        return ( dx + col * (dx + w))
      },
      y: function() {
        var h   = this.view.graph.node.height
        var dy  = this.view.graph.dy
        var row = this.node.Row

        return ( dy + row * (dy + h))
      },
      name: function() {
        return this.node.Name
      },
      type: function() {
        return this.node.Type
      }
    },
    methods: {
      // nodeSelected indicates that the node has been selected
      nodeSelected: function() {
        this.$emit('node-selected', this.node)
      }
    },
    template: `
      <g :transform="'translate(' + x + ',' + y + ')'">
        <foreignObject x=0 y=0 width=160 height=40>
          <div  class="node"
            :title="name"
            @click="nodeSelected">
            <div class="name">{{name}}</div>
            <div class="type">{{type}}</div>
          </div>
        </foreignObject>
        <node_destination
          v-bind:model="model"
          v-bind:view="view"
          v-bind:node="node"
          v-bind:destination="destination"
          v-bind:idx="idx"
          v-for="(destination,idx) in node.Destinations">
        </node_destination>
        <node_source
          v-bind:model="model"
          v-bind:view="view"
          v-bind:node="node"
          v-bind:source="source"
          v-bind:idx="idx"
          v-for="(source,idx) in node.Sources">
        </node_source>
      </g>`
  }
)

//------------------------------------------------------------------------------
