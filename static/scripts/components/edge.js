Vue.component(
  'edge',
  {
    props: ['model', 'view', 'edge'],
    computed: {
      path: function(){
        switch (this.edge.Type) {
          case "top-left":
            return topLeftPath(this.model, this.view, this.edge)
          case "above":
            return abovePath(this.model, this.view, this.edge)
          case "top-right":
            return topRightPath(this.model, this.view, this.edge)
          case "immediate-left":
            return immediateLeftPath(this.model, this.view, this.edge)
          case "immediate-below":
            return immediateBelowPath(this.model, this.view, this.edge)
          case "immediate-right":
            return immediateRightPath(this.model, this.view, this.edge)
          case "bottom-left":
            return bottomLeftPath(this.model, this.view, this.edge)
          case "below":
            return belowPath(this.model, this.view, this.edge)
          case "bottom-right":
            return bottomRightPath(this.model, this.view, this.edge)
        }

        return ""
      }
    },
    template: `
      <path class="relationship" :class="edge.Category" :d="path">
        <title>{{edge.Tag}}</title>
      </path>`
  }
)

//------------------------------------------------------------------------------

// topLeftPath creates a path from a source to a destination situated top left
function topLeftPath(model, view, edge) {
  var graph   = view.graph
  var layers  = model.Graph.Layers
  var columns = model.Graph.Columns
  var src     = edge.Source
  var dest    = edge.Destination
  var dx      = graph.dx
  var dy      = graph.dy
  var r       = graph.port.diameter/2
  var b       = graph.port.border
  var w       = graph.node.width
  var w2      = w - 2 * r
  var w3      = w + dx
  var h       = graph.node.height
  var h2      = h - 2 * r
  var h3      = h + dy
  var cx      = dx - 2 * r
  var cy      = dy - 4 * r

  var path = ""

  // start point below relationship port
  x1 = dx + edge.SrcCol * w3 + r + (edge.SrcIndex+1)/(src.Sources.length+1) * w2
  y1 = dy + edge.SrcRow * h3 + h + r

  path += "M " + x1 + " " + y1

  // line to source channel
  x2 = x1
  y2 = dy + edge.SrcRow * h3 + h + r + r + (edge.Channel1+1)/(layers[edge.SrcRow+1]+1) * cy - r

  path += "L " + x2 + " " + y2

  // curve right
  x3 = x2 - r
  y3 = y2
  x4 = x2 - r
  y4 = y2 + r

  path += "A " + x3 + " " + y3 + " 0 0 1 " + x4 + " " + y4

  // line along source channel
  x5 = dx + (edge.SrcCol-1) * w3 + r + (edge.Channel2+1)/(columns[edge.SrcCol]+1) * cx + r
  y5 = y4

  path += "L " + x5 + " " + y5

  // curve up
  x6 = x5
  y6 = y5 - r
  x7 = x5 - r
  y7 = y5 - r

  path += "A " + x6 + " " + y6 + " 0 0 1 " + x7 + " " + y7

  // line along vertical channel
  x8 = x7
  y8 = edge.DestRow * h3 + r + r + (edge.Channel3+1)/(layers[edge.DestRow]+1) * cy + r

  path += "L " + x8 + " " + y8

  // curve left
  x9  = x8 - r
  y9  = y8
  x10 = x8 - r
  y10 = y8 - r

  path += "A " + x6 + " " + y6 + " 0 0 0 " + x7 + " " + y7

  // line along destination channel
  x11 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2 + r
  y11 = y10

  path += "L " + x11 + " " + y11

  // curve down
  x12 = x11
  y12 = y11 + r
  x13 = x11 - r
  y13 = y11 + r

  path += "A " + x12 + " " + y12 + " 0 0 0 " + x13 + " " + y13

  // end point above destination port
  x14 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2
  y14 = dy + edge.DestRow * h3 - r

  path += "L " + x14 + " " + y14

  // finished
  return path
}

//------------------------------------------------------------------------------

// abovePath creates a path from a source to a destination situated above
function abovePath(model, view, edge) {
  var graph   = view.graph
  var layers  = model.Graph.Layers
  var columns = model.Graph.Columns
  var src     = edge.Source
  var dest    = edge.Destination
  var dx      = graph.dx
  var dy      = graph.dy
  var r       = graph.port.diameter/2
  var b       = graph.port.border
  var w       = graph.node.width
  var w2      = w - 2 * r
  var w3      = w + dx
  var h       = graph.node.height
  var h2      = h - 2 * r
  var h3      = h + dy
  var cx      = dx - 2 * r
  var cy      = dy - 4 * r

  var path = ""

  // start point below relationship port
  x1 = dx + edge.SrcCol * w3 + r + (edge.SrcIndex+1)/(src.Sources.length+1) * w2
  y1 = dy + edge.SrcRow * h3 + h + r

  path += "M " + x1 + " " + y1

  // line to source channel
  x2 = x1
  y2 = dy + edge.SrcRow * h3 + h + r + r + (edge.Channel1+1)/(layers[edge.SrcRow+1]+1) * cy - r

  path += "L " + x2 + " " + y2

  // curve left
  x3 = x2 + r
  y3 = y2
  x4 = x2 + r
  y4 = y2 + r

  path += "A " + x3 + " " + y3 + " 0 0 0 " + x4 + " " + y4

  // line along source channel
  x5 = dx + (edge.SrcCol+1) * w3 + r + (edge.Channel2+1)/(columns[edge.SrcCol]+1) * cx - r
  y5 = y4

  path += "L " + x5 + " " + y5

  // curve up
  x6 = x5
  y6 = y5 - r
  x7 = x5 + r
  y7 = y5 - r

  path += "A " + x6 + " " + y6 + " 0 0 0 " + x7 + " " + y7

  // line along vertical channel
  x8 = x7
  y8 = edge.DestRow * h3 + r + r + (edge.Channel3+1)/(layers[edge.DestRow]+1) * cy + r

  path += "L " + x8 + " " + y8

  // curve left
  x9  = x8 - r
  y9  = y8
  x10 = x8 - r
  y10 = y8 - r

  path += "A " + x6 + " " + y6 + " 0 0 0 " + x7 + " " + y7

  // line along destination channel
  x11 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2 + r
  y11 = y10

  path += "L " + x11 + " " + y11

  // curve down
  x12 = x11
  y12 = y11 + r
  x13 = x11 - r
  y13 = y11 + r

  path += "A " + x12 + " " + y12 + " 0 0 0 " + x13 + " " + y13

  // end point above destination port
  x14 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2
  y14 = dy + edge.DestRow * h3 - r

  path += "L " + x14 + " " + y14

  // finished
  return path
}

//------------------------------------------------------------------------------

// topRightPath creates a path from a source to a destination situated top right
function topRightPath(model, view, edge) {
  var graph   = view.graph
  var layers  = model.Graph.Layers
  var columns = model.Graph.Columns
  var src     = edge.Source
  var dest    = edge.Destination
  var dx      = graph.dx
  var dy      = graph.dy
  var r       = graph.port.diameter/2
  var b       = graph.port.border
  var w       = graph.node.width
  var w2      = w - 2 * r
  var w3      = w + dx
  var h       = graph.node.height
  var h2      = h - 2 * r
  var h3      = h + dy
  var cx      = dx - 2 * r
  var cy      = dy - 4 * r

  var path = ""

  // start point below relationship port
  x1 = dx + edge.SrcCol * w3 + r + (edge.SrcIndex+1)/(src.Sources.length+1) * w2
  y1 = dy + edge.SrcRow * h3 + h + r

  path += "M " + x1 + " " + y1

  // line to source channel
  x2 = x1
  y2 = dy + edge.SrcRow * h3 + h + r + r + (edge.Channel1+1)/(layers[edge.SrcRow+1]+1) * cy - r

  path += "L " + x2 + " " + y2

  // curve left
  x3 = x2 + r
  y3 = y2 + r

  path += "A " + r + " " + r + " 0 0 0 " + x3 + " " + y3

  // line along source channel
  x4 = (edge.SrcCol+1) * w3 + r + (edge.Channel2+1)/(columns[edge.SrcCol+1]+1) * cx - r
  y4 = y3

  path += "L " + x4 + " " + y4

  // curve up
  x5 = x4 + r
  y5 = y4 - r

  path += "A " + r + " " + r + " 0 0 0 " + x5 + " " + y5

  // line along vertical channel
  x6 = x5
  y6 = edge.DestRow * h3 + r + r + (edge.Channel3+1)/(layers[edge.DestRow]+1) * cy + r

  path += "L " + x6 + " " + y6

  // curve right
  x7 = x6 + r
  y7 = y6 - r

  path += "A " + r + " " + r + " 0 0 1 " + x7 + " " + y7

  // line along destination channel
  x8 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2 - r
  y8 = y7

  path += "L " + x8 + " " + y8

  // curve down
  x9 = x8 + r
  y9 = y8 + r

  path += "A " + r + " " + r + " 0 0 1 " + x9 + " " + y9

  // end point above destination port
  x10 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2
  y10 = dy + edge.DestRow * h3 - r

  path += "L " + x10 + " " + y10

  // finished
  return path
}

//------------------------------------------------------------------------------

// immediateLeftPath creates a path from a source to a destination situated immediate left
function immediateLeftPath(model, view, edge) {
  var graph   = view.graph
  var layers  = model.Graph.Layers
  var columns = model.Graph.Columns
  var src     = edge.Source
  var dest    = edge.Destination
  var dx      = graph.dx
  var dy      = graph.dy
  var r       = graph.port.diameter/2
  var b       = graph.port.border
  var w       = graph.node.width
  var w2      = w - 2 * r
  var w3      = w + dx
  var h       = graph.node.height
  var h2      = h - 2 * r
  var h3      = h + dy
  var cx      = dx - 2 * r
  var cy      = dy - 4 * r

  var path = ""

  // start point below relationship port
  x1 = dx + edge.SrcCol * w3 + r + (edge.SrcIndex+1)/(src.Sources.length+1) * w2
  y1 = dy + edge.SrcRow * h3 + h + r

  path += "M " + x1 + " " + y1

  // line to source channel
  x2 = x1
  y2 = dy + edge.SrcRow * h3 + h + r + r + (edge.Channel1+1)/(layers[edge.SrcRow+1]+1) * cy - r

  path += " L " + x2 + " " + y2

  // curve right
  x3 = x2 - r
  y3 = y2 + r

  path += " A " + r + " " + r + " 0 0 1 " + x3 + " " + y3

  // line along source channel
  x11 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2 + r
  y11 = y3

  path += " L " + x11 + " " + y11

  // curve down
  x13 = x11 - r
  y13 = y11 + r

  path += " A " + r + " " + r + " 0 0 0 " + x13 + " " + y13

  // end point above destination port
  x14 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2
  y14 = dy + edge.DestRow * h3 - r

  path += " L " + x14 + " " + y14

  // finished
  return path
}

//------------------------------------------------------------------------------

// immediateBelowPath creates a path from a source to a destination situated immediate below
function immediateBelowPath(model, view, edge) {
  var graph   = view.graph
  var layers  = model.Graph.Layers
  var columns = model.Graph.Columns
  var src     = edge.Source
  var dest    = edge.Destination
  var dx      = graph.dx
  var dy      = graph.dy
  var r       = graph.port.diameter/2
  var b       = graph.port.border
  var w       = graph.node.width
  var w2      = w - 2 * r
  var w3      = w + dx
  var h       = graph.node.height
  var h2      = h - 2 * r
  var h3      = h + dy
  var cx      = dx - 2 * r
  var cy      = dy - 4 * r

  var path = ""

  // relative points of source and destination port
  var p1 = (edge.SrcIndex+1)  / (src.Sources.length+1)
  var p2 = (edge.DestIndex+1) / (dest.Destinations.length+1)

  // start point below relationship port
  x1 = dx + edge.SrcCol * w3 + r + (edge.SrcIndex+1)/(src.Sources.length+1) * w2
  y1 = dy + edge.SrcRow * h3 + h + r

  path += "M " + x1 + " " + y1

  // if tight then draw a cubic spline
  if (Math.abs(p2 - p1) < 2 * r)  {
    x4 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2
    y4 = dy + edge.DestRow * h3 - r

    x2 = x1
    y2 = (y1 + y4) / 2
    x3 = x4
    y3 = y2

    path += "C " + x2 + "," + y2 + " " + x3 + "," + y3 + " " + x4 + "," + y4
  } else {
    // line to source channel
    x2 = x1
    y2 = dy + edge.SrcRow * h3 + h + r + r + (edge.Channel1+1)/(layers[edge.SrcRow+1]+1) * cy - r

    path += "L " + x2 + " " + y2

    // curve right or left
    if (p1 < p2) {
      // curve left
      x3 = x2 + r
      y3 = y2
      x4 = x2 + r
      y4 = y2 + r
      o  = 0
    } else {
      // curve right
      x3 = x2 - r
      y3 = y2
      x4 = x2 - r
      y4 = y2 + r
      o  = 1
    }

    path += "A " + x3 + " " + y3 + " 0 0 " + o + " " + x4 + " " + y4

    // line along source channel
    x11 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2 + r
    y11 = y10

    path += "L " + x11 + " " + y11

    // curve down
    if (p1 < p2) {
      x12 = x11
      y12 = y11 + r
      x13 = x11 + r
      y13 = y11 + r
      o   = 1
    } else {
      x12 = x11
      y12 = y11 + r
      x13 = x11 - r
      y13 = y11 + r
      o   = 1
    }

    path += "A " + x12 + " " + y12 + " 0 0 " + o + " " + x13 + " " + y13

    // end point above destination port
    x14 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2
    y14 = dy + edge.DestRow * h3 - r

    path += "L " + x14 + " " + y14
  }

  // finished
  return path
}

//------------------------------------------------------------------------------

// immediateRightPath creates a path from a source to a destination situated immediate right
function immediateRightPath(model, view, edge) {
  var graph   = view.graph
  var layers  = model.Graph.Layers
  var columns = model.Graph.Columns
  var src     = edge.Source
  var dest    = edge.Destination
  var dx      = graph.dx
  var dy      = graph.dy
  var r       = graph.port.diameter/2
  var b       = graph.port.border
  var w       = graph.node.width
  var w2      = w - 2 * r
  var w3      = w + dx
  var h       = graph.node.height
  var h2      = h - 2 * r
  var h3      = h + dy
  var cx      = dx - 2 * r
  var cy      = dy - 4 * r

  var path = ""

  // start point below relationship port
  x1 = dx + edge.SrcCol * w3 + r + (edge.SrcIndex+1)/(src.Sources.length+1) * w2
  y1 = dy + edge.SrcRow * h3 + h + r

  path += "M " + x1 + " " + y1

  // line to source channel
  x2 = x1
  y2 = dy + edge.SrcRow * h3 + h + r + r + (edge.Channel1+1)/(layers[edge.SrcRow+1]+1) * cy - r

  path += "L " + x2 + " " + y2

  // curve left
  x3 = x2 + r
  y3 = y2 + r

  path += "A " + r + " " + r + " 0 0 0 " + x3 + " " + y3

  // line along source channel
  x4 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2 - r
  y4 = y3

  path += "L " + x4 + " " + y4

  // curve down
  x5 = x4 + r
  y5 = y4 + r

  path += "A " + r + " " + r + " 0 0 1 " + x5 + " " + y5

  // end point above destination port
  x6 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2
  y6 = dy + edge.DestRow * h3 - r

  path += "L " + x6 + " " + y6

  // finished
  return path
}

//------------------------------------------------------------------------------

// bottomLeftPath creates a path from a source to a destination situated bottom left
function bottomLeftPath(model, view, edge) {
  var graph   = view.graph
  var layers  = model.Graph.Layers
  var columns = model.Graph.Columns
  var src     = edge.Source
  var dest    = edge.Destination
  var dx      = graph.dx
  var dy      = graph.dy
  var r       = graph.port.diameter/2
  var b       = graph.port.border
  var w       = graph.node.width
  var w2      = w - 2 * r
  var w3      = w + dx
  var h       = graph.node.height
  var h2      = h - 2 * r
  var h3      = h + dy
  var cx      = dx - 2 * r
  var cy      = dy - 4 * r

  var path = ""

  // start point below relationship port
  x1 = dx + edge.SrcCol * w3 + r + (edge.SrcIndex+1)/(src.Sources.length+1) * w2
  y1 = dy + edge.SrcRow * h3 + h + r

  path += "M " + x1 + " " + y1

  // line to source channel
  x2 = x1
  y2 = dy + edge.SrcRow * h3 + h + r + r + (edge.Channel1+1)/(layers[edge.SrcRow+1]+1) * cy - r

  path += "L " + x2 + " " + y2

  // curve right
  x3 = x2 - r
  y3 = y2 + r

  path += "A " + r + " " + r + " 0 0 1 " + x3 + " " + y3

  // line along source channel
  x4 = edge.SrcCol * w3 + r + (edge.Channel2+1)/(columns[edge.SrcCol]+1) * cx - r
  y4 = y3

  path += "L " + x4 + " " + y4

  // curve down
  x5 = x4 - r
  y5 = y4 + r

  path += "A " + r + " " + r + " 0 0 0 " + x5 + " " + y5

  // line along vertical channel
  x6 = x5
  y6 = edge.DestRow * h3 + r + r + (edge.Channel3+1)/(layers[edge.DestRow]+1) * cy + r

  path += "L " + x6 + " " + y6

  // curve right
  x7 = x6 - r
  y7 = y6 + r

  path += "A " + r + " " + r + " 0 0 1 " + x7 + " " + y7

  // line along destination channel
  x8 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2 + r
  y8 = y7

  path += "L " + x8 + " " + y8

  // curve down
  x9 = x8 - r
  y9 = y8 + r

  path += "A " + r + " " + r + " 0 0 0 " + x9 + " " + y9

  // end point above destination port
  x10 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2
  y10 = dy + edge.DestRow * h3 - r

  path += "L " + x10 + " " + y10

  // finished
  return path
}

//------------------------------------------------------------------------------

// belowPath creates a path from a source to a destination situated below
function belowPath(model, view, edge) {
  var graph   = view.graph
  var layers  = model.Graph.Layers
  var columns = model.Graph.Columns
  var src     = edge.Source
  var dest    = edge.Destination
  var dx      = graph.dx
  var dy      = graph.dy
  var r       = graph.port.diameter/2
  var b       = graph.port.border
  var w       = graph.node.width
  var w2      = w - 2 * r
  var w3      = w + dx
  var h       = graph.node.height
  var h2      = h - 2 * r
  var h3      = h + dy
  var cx      = dx - 2 * r
  var cy      = dy - 4 * r

  var path = ""

  // start point below relationship port
  x1 = dx + edge.SrcCol * w3 + r + (edge.SrcIndex+1)/(src.Sources.length+1) * w2
  y1 = dy + edge.SrcRow * h3 + h + r

  path += "M " + x1 + " " + y1

  // line to source channel
  x2 = x1
  y2 = dy + edge.SrcRow * h3 + h + r + r + (edge.Channel1+1)/(layers[edge.SrcRow+1]+1) * cy - r

  path += " L " + x2 + " " + y2

  // curve left
  x3 = x2 + r
  y3 = y2 + r

  path += " A " + r + " " + r + " 0 0 0 " + x3 + " " + y3

  // line along source channel
  x4 = (edge.SrcCol+1) * w3 + r + (edge.Channel2+1)/(columns[edge.SrcCol]+1) * cx - r
  y4 = y3

  path += " L " + x4 + " " + y4

  // curve down
  x5 = x4 + r
  y5 = y4 + r

  path += "A " + r + " " + r + " 0 0 1 " + x5 + " " + y5

  // line along vertical channel
  x6 = x5
  y6 = edge.DestRow * h3 + r + r + (edge.Channel3+1)/(layers[edge.DestRow]+1) * cy + r

  path += " L " + x6 + " " + y6

  // curve right
  x7 = x6 - r
  y7 = y6 + r

  path += "A " + r + " " + r + " 0 0 1 " + x7 + " " + y7

  // line along destination channel
  x8 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2 + r
  y8 = y7

  path += " L " + x8 + " " + y8

  // curve down
  x9 = x8 - r
  y9 = y8 + r

  path += " A " + r + " " + r + " 0 0 0 " + x9 + " " + y9

  // end point above destination port
  x10 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2
  y10 = dy + edge.DestRow * h3 - r

  path += " L " + x10 + " " + y10

  // finished
  return path
}

//------------------------------------------------------------------------------

// bottomRightPath creates a path from a source to a destination situated bottom right
function bottomRightPath(model, view, edge) {
  var graph   = view.graph
  var layers  = model.Graph.Layers
  var columns = model.Graph.Columns
  var src     = edge.Source
  var dest    = edge.Destination
  var dx      = graph.dx
  var dy      = graph.dy
  var r       = graph.port.diameter/2
  var b       = graph.port.border
  var w       = graph.node.width
  var w2      = w - 2 * r
  var w3      = w + dx
  var h       = graph.node.height
  var h2      = h - 2 * r
  var h3      = h + dy
  var cx      = dx - 2 * r
  var cy      = dy - 4 * r

  var path = ""

  // start point below relationship port
  x1 = dx + edge.SrcCol * w3 + r + (edge.SrcIndex+1)/(src.Sources.length+1) * w2
  y1 = dy + edge.SrcRow * h3 + h + r

  path += "M " + x1 + " " + y1

  // line to source channel
  x2 = x1
  y2 = dy + edge.SrcRow * h3 + h + r + r + (edge.Channel1+1)/(layers[edge.SrcRow+1]+1) * cy - r

  path += "L " + x2 + " " + y2

  // curve left
  x3 = x2 + r
  y3 = y2
  x4 = x2 + r
  y4 = y2 + r

  path += "A " + x3 + " " + y3 + " 0 0 0 " + x4 + " " + y4

  // line along source channel
  x5 = dx + (edge.SrcCol+1) * w3 + r + (edge.Channel2+1)/(columns[edge.SrcCol+1]+1) * cx + r
  y5 = y4

  path += "L " + x5 + " " + y5

  // curve down
  x6 = x5
  y6 = y5 + r
  x7 = x5 + r
  y7 = y5 + r

  path += "A " + x6 + " " + y6 + " 0 0 1 " + x7 + " " + y7

  // line along vertical channel
  x8 = x7
  y8 = edge.DestRow * h3 + r + r + (edge.Channel3+1)/(layers[edge.DestRow]+1) * cy + r

  path += "L " + x8 + " " + y8

  // curve left
  x9  = x8 + r
  y9  = y8
  x10 = x8 + r
  y10 = y8 + r

  path += "A " + x6 + " " + y6 + " 0 0 0 " + x7 + " " + y7

  // line along destination channel
  x11 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2 + r
  y11 = y10

  path += "L " + x11 + " " + y11

  // curve down
  x12 = x11
  y12 = y11 + r
  x13 = x11 + r
  y13 = y11 + r

  path += "A " + x12 + " " + y12 + " 0 0 1 " + x13 + " " + y13

  // end point above destination port
  x14 = dx + edge.DestCol * w3 + r + (edge.DestIndex+1)/(dest.Destinations.length+1) * w2
  y14 = dy + edge.DestRow * h3 - r

  path += "L " + x14 + " " + y14

  // finished
  return path
}

//------------------------------------------------------------------------------
