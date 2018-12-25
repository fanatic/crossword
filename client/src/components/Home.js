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
    let game_id = '';
    if (props.game) {
      game_id = props.game.id;
    }
    this.state = { name: props.player_name, game_id: game_id };
  }

  componentWillReceiveProps(newProps) {
    if (newProps.player_name !== this.props.player_name) {
      this.setState({ name: newProps.player_name });
    }
    if (newProps.game !== this.props.game) {
      let game_id = '';
      if (newProps.game) {
        game_id = newProps.game.id;
      }
      this.setState({ game_id });
    }
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

    this.props.postGame(result => {
      this.props.setGameID(result.id, '');
      return {};
    });
  };

  handleJoin = event => {
    event.preventDefault();

    this.props.postPlayer(this.state.game_id, { name: this.state.name }, result => {
      this.props.setGameID(result.id, this.state.name);
      return {};
    });
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
  postPlayer: (game_id, body, callback) => ({
    postPlayerResponse: {
      url: `http://localhost:8080/games/${game_id}/players`,
      method: 'POST',
      body: JSON.stringify(body),
      andThen: callback
    }
  }),
  postGame: callback => ({
    postGameResponse: {
      url: `http://localhost:8080/games`,
      method: 'POST',
      andThen: callback
    }
  })
}))(Home);
