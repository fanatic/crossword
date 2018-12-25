import React, { Component } from 'react';
import { connect, PromiseState } from 'react-refetch';
import PropTypes from 'prop-types';

class Home extends Component {
  static propTypes = {
    postPlayer: PropTypes.func.isRequired,
    postPlayerResponse: PropTypes.instanceOf(PromiseState),
    postGame: PropTypes.func.isRequired,
    postGameResponse: PropTypes.instanceOf(PromiseState)
  };

  constructor(props) {
    super(props);
    this.state = { name: '', game_id: '' };
  }

  handleInputChange = event => {
    const target = event.target;
    const value = target.type === 'checkbox' ? target.checked : target.value;
    const name = target.name;

    this.setState({
      [name]: value
    });
  };

  handleHost = event => {
    event.preventDefault();

    this.props.postGame();
  };

  handleJoin = event => {
    event.preventDefault();

    this.props.postPlayer(this.state.game_id, { name: this.state.name });
  };

  render() {
    return (
      <div>
        <h2>Home Screen</h2>
        <form onSubmit={this.handleHost}>
          <input type="submit" value="Host Game" />
        </form>
        <form onSubmit={this.handleJoin}>
          <input type="text" placeholder="Name" name="name" value={this.state.name} onChange={this.handleInputChange} />
          <input
            type="text"
            placeholder="Game ID"
            name="game_id"
            value={this.state.game_id}
            onChange={this.handleInputChange}
          />
          <input type="submit" value="Join Game" />
        </form>
      </div>
    );
  }
}

export default connect(props => ({
  postPlayer: (game_id, body) => ({
    postPlayerResponse: {
      url: `http://localhost:8080/games/${game_id}/players`,
      method: 'POST',
      body: JSON.stringify(body)
    }
  }),
  postGame: () => ({
    postGameResponse: {
      url: `http://localhost:8080/games`,
      method: 'POST'
    }
  })
}))(Home);
