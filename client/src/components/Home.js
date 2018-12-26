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
    this.state = { name: props.player_name, game_id: game_id, board_id: '', board_url: '' };
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

    this.props.postGame({ board_id: this.state.board_id }, result => {
      this.props.postPlayer(result.id, { name: this.state.name }, result => {
        this.props.setGameID(result.id, this.state.name);
        return {};
      });
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

  handleSubmitBoard = event => {
    event.preventDefault();

    this.props.postLayout({ url: this.state.board_url }, result => {
      this.setState({ board_url: '' });
      return {};
    });
  };

  render() {
    const { fetchLayouts } = this.props;

    return (
      <div>
        <h2>Home Screen</h2>
        <form onSubmit={this.handleHost}>
          <input type="text" placeholder="Name" name="name" value={this.state.name} onChange={this.handleInputChange} />
          <select name="board_id" value={this.state.board_id} onChange={this.handleInputChange}>
            <option value="" />
            {fetchLayouts.fulfilled &&
              fetchLayouts.value.map(l => (
                <option key={l.id} value={l.id}>
                  {l.id}
                </option>
              ))}
          </select>
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
        <form onSubmit={this.handleSubmitBoard}>
          <input
            type="text"
            placeholder="Board URL"
            name="board_url"
            value={this.state.board_url}
            onChange={this.handleInputChange}
          />
          <input type="submit" value="Upload Board" />
        </form>
      </div>
    );
  }
}

export default connect(props => ({
  postPlayer: (game_id, body, callback) => ({
    postPlayerResponse: {
      url: `http://192.168.3.38:8080/games/${game_id}/players`,
      method: 'POST',
      body: JSON.stringify(body),
      andThen: callback
    }
  }),
  postGame: (body, callback) => ({
    postGameResponse: {
      url: `http://192.168.3.38:8080/games`,
      method: 'POST',
      body: JSON.stringify(body),
      andThen: callback
    }
  }),
  fetchLayouts: {
    url: `http://192.168.3.38:8080/layouts`
  },
  postLayout: (body, callback) => ({
    fetchLayouts: {
      url: `http://192.168.3.38:8080/layouts`,
      method: 'POST',
      body: JSON.stringify(body),
      andThen: callback
    }
  })
}))(Home);
