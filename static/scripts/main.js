var app

function main() {
  // load data
  loadModel()

  // present user interface
  app = new Vue({
    el:   '#app',
    data: {
      model: model,
      view: view,
    },
    template: `<app v-bind:model="model" v-bind:view="view"></app>`
  })
}


window.onload = main;
