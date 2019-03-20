// SolutionGraph represents a graph derived from a solution
function SolutionGraph(model, view, domain, solution) {
  // member variables
  this.model         = model       // model to hold the solution graph
  this.view          = view        // view containing the dimension parameters
  this.domain        = domain      // domain of the solution
  this.solution      = solution    // name of the solution
  this.Components    = {}          // components of the catalog
  this.Elements      = {}          // elements of the solution
  this.Clusters      = {}          // clusters of the elements
  this.Relationships = {}          // relationships between clusters
  this.Nodes         = {}          // nodes of the graph
  this.Sources       = []          // source ports of the nodes
  this.Destinations  = []          // destination ports of the nodes
  this.Edges         = {}          // edges of the graph
  this.Layers        = []          // layer channels of the graph
  this.Columns       = []          // column channels of the graph
  this.Width         = 0           // width of the graph
  this.Height        = 0           // height of the graph

  //
  // load catalog, architecture and solution
  loadAll(domain, solution)
  // initialise lookups of graph
  .then(() => {
    this.calculateComponents()
    this.calculateElements()
    this.calculateClusters()
    this.calculateRelationships()
    this.calculateNodes()
    this.calculateSources()
    this.calculateDestinations()
    this.sortSources()
    this.sortDestinations()
    this.sortNodes()
    this.calculateEdges()
    this.calculateDimensions()
    this.model.Graph = this
  })
  .catch((error) => {
    console.log(error);
  })
}

//------------------------------------------------------------------------------

// calculateComponents determines all components
SolutionGraph.prototype.calculateComponents = function() {
  var catalog    = this.model.Catalog
  var components = this.Components

  Object.values(catalog).forEach((component) => {
    var name = component.Component + " - " + component.Version

    components[name] = component
  })
}

//------------------------------------------------------------------------------

// calculateElements determines all solution elements
SolutionGraph.prototype.calculateElements = function() {
  var elements1 = this.model.Solution.Elements
  var elements2 = this.Elements

  Object.values(elements1).forEach((element) => {
    var name = element.Element

    elements2[name] = element
  })
}

//------------------------------------------------------------------------------

// calculateClusters determines all clusters of the solution elements
SolutionGraph.prototype.calculateClusters = function() {
  var elements = this.Elements
  var clusters = this.Clusters

  Object.values(elements).forEach((element) => {
    Object.values(element.Clusters).forEach((cluster) => {
      var name = element.Element + " / " + cluster.Version

      clusters[name] = cluster
    })
  })
}

//------------------------------------------------------------------------------

// calculateRelationships determines all relationships of the clusters of the solution elements
SolutionGraph.prototype.calculateRelationships = function() {
  var elements      = this.Elements
  var relationships = this.Relationships

  Object.values(elements).forEach((element) => {
    Object.values(element.Clusters).forEach((cluster) => {
      Object.values(cluster.Relationships).forEach((relationship) => {
        var name = element.Element + " / " + cluster.Version + " / " + relationship.Relationship

        relationships[name] = relationship
      })
    })
  })
}

//------------------------------------------------------------------------------

// calculateNodes determines all nodes of the graph
SolutionGraph.prototype.calculateNodes = function() {
  var elements = this.Elements
  var nodes    = this.Nodes

  Object.values(elements).forEach((element) => {
    var name = element.Element

    nodes[name] = {
      Name:          element.Element,
      Type:          element.Component,
      Inbound1:      0,        // number of inbound context relationships
      Inbound2:      0,        // number of inbound service relationships
      Row:           1000000,  // initialise with a very large number
      Column:        1000000,  // initialise with a very large number
      Sources:       [],
      Destinations:  [],
      Element:       element
    }
  })
}

//------------------------------------------------------------------------------

// calculateDestinations determines all service ports of the nodes
SolutionGraph.prototype.calculateDestinations = function() {
  var elements     = this.Elements
  var nodes        = this.Nodes
  var destinations = this.Destinations

  Object.values(elements).forEach((element) => {
    Object.values(element.Clusters).forEach((cluster) => {
      var name   = element.Element + " / " + cluster.Version
      var tag    = cluster.Version
      var state  = cluster.State
      var target = cluster.Target
      var node   = nodes[element.Element]

      // create a new destination
      var destination = {
        Name:    name,
        Tag:     tag,
        State:   state,
        Target:  target,
        Node:    node,
        Index:   -1
      }

      // add destinations to the graph and node
      destinations[name] = destination
      node.Destinations.push(destination)
    })
  })
}

//------------------------------------------------------------------------------

// calculateSources determines all relationships of the nodes
SolutionGraph.prototype.calculateSources = function() {
  var elements = this.Elements
  var nodes    = this.Nodes
  var sources  = this.Sources

  Object.values(elements).forEach((element) => {
    Object.values(element.Clusters).forEach((cluster) => {
      Object.values(cluster.Relationships).forEach((relationship) => {
        var name = element.Element + " / " + cluster.Version + " / " + relationship.Relationship
        var tag  = relationship.Relationship
        var node = nodes[element.Element]

        // create a new source
        var source = {
          Name:         name,
          Tag:          tag,
          Node:         node,
          Index:        -1,
          Relationship: relationship
        }

        // add source to the graph and node
        sources[name] = source
        node.Sources.push(source)

        // determine target and increment their inbound counters
        var targetNode = nodes[relationship.Element]

        if (relationship.Type == "context") {
          targetNode.Inbound1++
        } else {
          targetNode.Inbound2++
        }
      })
    })
  })
}

//------------------------------------------------------------------------------

// sortSources sort the sources of a node
SolutionGraph.prototype.sortSources = function() {
  var nodes = this.Nodes

  // sort the sources
  Object.values(nodes).forEach((node) => {
    node.Sources.sort((a,b) => {return a.Tag < b.Tag ? -1 : a.Tag > b.Tag ? 1 : 0})
  })

  // update the index
  Object.values(nodes).forEach((node) => {
    Object.values(node.Sources).forEach((source,index) => {
      source.Index = index
    })
  })
}

//------------------------------------------------------------------------------

// sortDestinations sort the destinations of a node
SolutionGraph.prototype.sortDestinations = function() {
  var nodes = this.Nodes

  // sort the destinations
  Object.values(nodes).forEach((node) => {
    node.Destinations.sort((a,b) => {return a.Tag < b.Tag ? -1 : a.Tag > b.Tag ? 1 : 0})
  })

  // update the index
  Object.values(nodes).forEach((node) => {
    Object.values(node.Destinations).forEach((destination,index) => {
      destination.Index = index
    })
  })
}

//------------------------------------------------------------------------------

// sortNodes sort the nodes of the graph (Kahn's Algorithm)
SolutionGraph.prototype.sortNodes = function() {
  var nodes  = this.Nodes
  var list   = []
  var queue1 = []
  var queue2 = []
  var node   = null

  // collect all nodes with no inbound connections
  Object.values(nodes).forEach((node) => {
    if (node.Inbound1 == 0) {
      node.Row = 0
      queue1.push(node)
    }
    if (node.Inbound2 == 0) {
      node.Column = 0
      queue2.push(node)
    }
  })

  // iterate until queue 1 is empty
  while (queue1.length > 0) {
    // add next node in queue to list
    node = queue1.pop()

    // determine relationship targets and decrement their inbound counters
    Object.values(node.Sources).forEach((source) => {
      var relationship = source.Relationship
      var targetNode   = nodes[relationship.Element]

      if (relationship.Type == "context") {
        targetNode.Inbound1--

        // add target node to queue if inbound counter has reached 0
        if (targetNode.Inbound1 == 0) {
          targetNode.Row = node.Row + 1
          queue1.push(targetNode)
        }
      }
    })
  }

  // iterate until queue 2 is empty
  while (queue2.length > 0) {
    // add next node in queue to list
    node = queue2.pop()

    // determine relationship targets and decrement their inbound counters
    Object.values(node.Sources).forEach((source) => {
      var relationship = source.Relationship
      var targetNode   = nodes[relationship.Element]

      if (relationship.Type == "service") {
        targetNode.Inbound2--

        // add target node to queue if inbound counter has reached 0
        if (targetNode.Inbound2 == 0) {
          targetNode.Column = node.Column + 1
          queue2.push(targetNode)
        }
      }
    })
  }

  // sort columns
  Object.values(nodes).forEach((node) => { list.push(node) })
  list.sort((a,b) => {
    if      (a.Row    < b.Row   ) { return -1 }
    else if (a.Row    > b.Row   ) { return  1 }
    else if (a.Column < b.Column) { return -1 }
    else if (a.Column > b.Column) { return  1 }
    return 0
  })

  // adjust column index of nodes
  layer     = 0
  column    = 0
  maxLayer  = 0
  maxColumn = 0
  Object.values(list).forEach((node) => {
    maxLayer = Math.max(maxLayer, node.Row+1)
    if (node.Row == layer) {
      node.Column = column
      column++
      maxColumn = Math.max(maxColumn, column)
    } else {
      layer       = node.Row
      column      = 0
      node.Column = column
      column++
      maxColumn = Math.max(maxColumn, column)
    }
  })

  // initialise layers and columns paths
  // these indicate how many relationships pass through this path
  for (var i = 0; i < maxLayer; i++) {
    this.Layers.push(0)
  }
  for (var i = 0; i <= maxColumn; i++) {
    this.Columns.push(0)
  }
}

//------------------------------------------------------------------------------

// calculateEdges determines the required paths for the edges
SolutionGraph.prototype.calculateEdges = function() {
  var nodes        = this.Nodes
  var edges        = this.Edges
  var destinations = this.Destinations
  var layers       = this.Layers
  var columns      = this.Columns

  Object.values(nodes).forEach((node) => {
    Object.values(node.Sources).forEach((source) => {
      var relationship = source.Relationship
      var destination  = destinations[relationship.Element + " / " + relationship.Version]
      var name         = source.Name + " / " + destination.Name
      var tag          = source.Name + " --> " + destination.Name
      var category     = relationship.Type

      var edge = {
        Tag:         tag,
        Source:      source.Node,
        Destination: destination.Node,
        SrcRow:      source.Node.Row,
        SrcCol:      source.Node.Column,
        SrcIndex:    source.Index,
        DestRow:     destination.Node.Row,
        DestCol:     destination.Node.Column,
        DestIndex:   destination.Index,
        Category:    category,
        Type:        "",
        Channel1:    -1,
        Channel2:    -1,
        Channel3:    -1
      }

      // case A: dest to the top left of src
      if (edge.DestRow <= edge.SrcRow && edge.DestCol < edge.SrcCol) {
        edge.Type     = "top-left"
        edge.Channel1 = layers[edge.SrcRow+1]++
        edge.Channel2 = columns[edge.SrcCol]++
        edge.Channel3 = layers[edge.DestRow]++
      }

      // case B: dest above of src
      if (edge.DestRow <= edge.SrcRow && edge.DestCol == edge.SrcCol) {
        edge.Type     = "above"
        edge.Channel1 = layers[edge.SrcRow+1]++
        edge.Channel2 = columns[edge.SrcCol+1]++
        edge.Channel3 = layers[edge.DestRow]++
      }

      // case C: dest to the top right of src
      if (edge.DestRow <= edge.SrcRow && edge.SrcCol < edge.DestCol) {
        edge.Type     = "top-right"
        edge.Channel1 = layers[edge.SrcRow+1]++
        edge.Channel2 = columns[edge.SrcCol+1]++
        edge.Channel3 = layers[edge.DestRow]++
      }

      // case D: dest to the immediate left of src
      if ((edge.SrcRow+1) == edge.DestRow && edge.DestCol < edge.SrcCol) {
        edge.Type     = "immediate-left"
        edge.Channel1 = layers[edge.SrcRow+1]++
      }

      // case E: dest immediately below of src
      if ((edge.SrcRow+1) == edge.DestRow && edge.DestCol == edge.SrcCol) {
        edge.Type     = "immediate-below"
        edge.Channel1 = layers[edge.SrcRow+1]++
      }

      // case F: dest to the immediate right of src
      if ((edge.SrcRow+1) == edge.DestRow && edge.SrcCol < edge.DestCol) {
        edge.Type     = "immediate-right"
        edge.Channel1 = layers[edge.SrcRow+1]++
      }

      // case G: dest to the bottom left of src
      if ((edge.SrcRow+1) < edge.DestRow && edge.DestCol < edge.SrcCol) {
        edge.Type     = "bottom-left"
        edge.Channel1 = layers[edge.SrcRow+1]++
        edge.Channel2 = columns[edge.SrcCol]++
        edge.Channel3 = layers[edge.DestRow]++
      }

      // case H: dest below of src
      if ((edge.SrcRow+1) < edge.DestRow && edge.DestCol == edge.SrcCol) {
        edge.Type     = "below"
        edge.Channel1 = layers[edge.SrcRow+1]++
        edge.Channel2 = columns[edge.SrcCol+1]++
        edge.Channel3 = layers[edge.DestRow]++
      }

      // case I: dest to the bottom right of src
      if ((edge.SrcRow+1) < edge.DestRow && edge.SrcCol < edge.DestCol ) {
        edge.Type     = "bottom-right"
        edge.Channel1 = layers[edge.SrcRow+1]++
        edge.Channel2 = columns[edge.SrcCol+1]++
        edge.Channel3 = layers[edge.DestRow]++
      }

      // add egde
      edges[name] = edge
    })
  })
}

//------------------------------------------------------------------------------

// calculateDimensions determines the width and height of the graph
SolutionGraph.prototype.calculateDimensions = function() {
  vg = this.view.graph

  this.Width  = vg.dx + (this.Columns.length-1) * (vg.node.width  + vg.dx)
  this.Height = vg.dy + this.Layers.length      * (vg.node.height + vg.dy)
}

//------------------------------------------------------------------------------
