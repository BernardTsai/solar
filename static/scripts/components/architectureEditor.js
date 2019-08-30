Vue.component(
  'architectureEditor',
  {
    props: ['model','view'],
    template: `
      <div class="architectureEditor">
        <div class="header">
          <h3>Architecture: {{model.Architecture.Architecture}} - {{model.Architecture.Version}}</h3>
        </div>

        <table style="width: 100%">
          <col width="10*">
          <col width="990*">
          <tr>
            <td><strong>Architecture:</strong></td>
            <td>
              <input type="text" readonly v-model="model.Architecture.Architecture"/>
            </td>
          </tr>
          <tr>
            <td>&nbsp;Version:</td>
            <td>
              <input type="text" readonly v-model="model.Architecture.Version"/>
            </td>
          </tr>
          <tr>
            <td>&nbsp;Conf.&nbsp;Parameters:</td>
            <td>
              <textarea id="configuration" rows=20
                v-model="model.Architecture.Configuration">
              </textarea>
            </td>
          </tr>
        </table>


      </div>`
  }
)
