import ReactDOM from 'react-dom/client';
import App from './App';
import {BrowserRouter} from "react-router-dom";
import {Provider} from "react-redux";
import Store from './store'
import { onAuthStateChanged } from 'firebase/auth';
import { restoreSession } from './model/user';
import { auth } from './auth';


const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
      <BrowserRouter>
          <Provider store={Store}>
                <App />
          </Provider>
      </BrowserRouter>
);

onAuthStateChanged(auth, user => {
  Store.dispatch(restoreSession(user))
})
