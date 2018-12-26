import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import GameWrapper from './components/GameWrapper';
import TVWrapper from './components/TVWrapper';

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
          <Route exact path="/" component={GameWrapper} />
          <Route exact path="/tv/:id" component={TVWrapper} />
        </div>
      </Router>
    );
  }
}

export default App;
