import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import GameWrapper from './components/GameWrapper';

import './App.css';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      game_id: localStorage.getItem('game_id')
    };
  }

  setGameID = game_id => {
    this.setState({ game_id });
    localStorage.setItem('game_id', game_id);
  };

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

          <Route exact path="/" component={GameWrapper} />
        </div>
      </Router>
    );
  }
}

export default App;
