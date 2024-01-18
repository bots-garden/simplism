import { LitElement, html, css} from 'lit-element';
import { myStyles as myStyles } from './styles.js';

class MyTitle extends LitElement {
  static styles = [myStyles]

  render() {
    return html`
    <h1 class="title"> {{message}} </h1>
    `;
  }
}

customElements.define('my-title', MyTitle);

class MainApp extends LitElement {
  static styles = [myStyles]

  render() {
    return html`
    <section class="container">
      <div>
        <my-title></my-title>
      </div>
    </section>    
    `;
  }
}

customElements.define('main-app', MainApp);
