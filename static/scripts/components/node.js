Vue.component(
  'node_destination',
  {
    props: ['model', 'view', 'node', 'destination', 'idx'],
    computed: {
      left: function() {
        var portIndex = parseInt(this.idx) + 1
        var ports     = this.node.Destinations.length + 1
        var radius    = this.view.graph.port.diameter/2
        var w         = this.view.graph.node.width - 2 * radius

        return (radius + portIndex/ports*w - radius) + 'px'
      },
      tag: function() {
        return this.destination.Tag
      }
    },
    template: `
      <div class="destination" :class="destination.State"
        v-bind:style="{left: left}"
        v-bind:title="tag">
      </div>`
  }
)

//------------------------------------------------------------------------------

Vue.component(
  'node_source',
  {
    props: ['model', 'view', 'node', 'source', 'idx'],
    computed: {
      left: function() {
        var portIndex = parseInt(this.idx) + 1
        var ports     = this.node.Sources.length + 1
        var radius    = this.view.graph.port.diameter/2
        var w         = this.view.graph.node.width - 2 * radius

        return (radius + portIndex/ports*w - radius) + 'px'
      },
      tag: function() {
        return this.source.Tag
      }
    },
    template: `
      <div class="source"
        v-bind:style="{left: left}"
        v-bind:title="tag">
      </div>`
  }
)

//------------------------------------------------------------------------------

Vue.component(
  'node',
  {
    props: ['model', 'view', 'node'],
    computed: {
      left: function() {
        var w   = this.view.graph.node.width
        var dx  = this.view.graph.dx
        var col = this.node.Column

        return ( dx + col * (dx + w)) + "px"
      },
      top: function() {
        var h   = this.view.graph.node.height
        var dy  = this.view.graph.dy
        var row = this.node.Row

        return ( dy + row * (dy + h)) + "px"
      },
      name: function() {
        return this.node.Name
      },
      type: function() {
        return this.node.Type
      }
    },
    template: `
      <div  class="node"
        v-bind:title="name"
        v-bind:style="{left: left, top: top}">
        <div class="name">{{name}}</div>
        <div class="type">{{type}}</div>
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
      </div>`
  }
)

//------------------------------------------------------------------------------
