import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import Debug from './components/Debug';

import './App.css';

class App extends Component {
  render() {
    return (
      <Router>
        <div>
          <h1>Crossword</h1>
          {/* <ul>
            <li>
              <Link to="/">Home</Link>
            </li>
            <li>
              <Link to="/about">About</Link>
            </li>
            <li>
              <Link to="/topics">Topics</Link>
            </li>
          </ul>

          <hr /> */}

          <Route exact path="/" component={Debug} />
        </div>
      </Router>
    );
  }
}

export default App;
