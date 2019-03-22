Vue.component(
  'graph',
  {
    props: ['model', 'view', 'graph'],
    methods: {
      // edgeSelected indicates that an edge has been selected
      edgeSelected: function(edge) {
        this.$emit('edge-selected', edge)
      },
      // nodeSelected indicates that a node has been selected
      nodeSelected: function(node) {
        this.$emit('node-selected', node)
      }
    },
    template: `
      <svg id="canvas" v-bind:style="{ width: graph.Width + 'px', height: graph.Height + 'px'}">
        <edge :model="model" :view="view" :edge="edge" v-for="edge in graph.Edges" @edge-selected="edgeSelected"/>
        <node :model="model" :view="view" :node="node" v-for="node in graph.Nodes" @node-selected="nodeSelected"/>
      </svg>`
  }
)
