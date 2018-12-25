import React, { Component } from 'react';
import { connect, PromiseState } from 'react-refetch';
import PropTypes from 'prop-types';
import BoardFragment from './BoardFragment';

class Mobile extends Component {
  static propTypes = {
    game: PropTypes.object.isRequired,
    player_name: PropTypes.string.isRequired,
    postGuess: PropTypes.func.isRequired,
    postGuessResponse: PropTypes.instanceOf(PromiseState)
  };

  constructor(props) {
    super(props);
    this.state = { guess: '' };
  }

  handleInputChange = event => {
    const target = event.target;
    const value = target.type === 'checkbox' ? target.checked : target.value;
    const name = target.name;

    this.setState({
      [name]: value
    });
  };

  handleGuess = event => {
    event.preventDefault();

    this.props.postGuess(
      this.props.game.id,
      { player_name: this.props.player_name, guess: this.state.guess },
      result => {
        this.props.setGameID(result.id);
        return {};
      }
    );
  };

  render() {
    const { player_name, game } = this.props;

    return (
      <div>
        <h2>
          Mobile Play ({player_name} - {game.id})
        </h2>
        {game.current_clue && (
          <React.Fragment>
            <p>{game.current_clue.description}</p>
            <BoardFragment clue={game.current_clue} />
            <form onSubmit={this.handleGuess}>
              <input type="text" name="guess" value={this.state.guess} onChange={this.handleInputChange} />
              <input type="submit" value="Guess" />
            </form>
            <br />
            <small>Waiting on {game.current_players.map(p => p.name).join(', ')}.</small>
          </React.Fragment>
        )}
      </div>
    );
  }
}

export default connect(props => ({
  postGuess: (game_id, body, callback) => ({
    postGuessResponse: {
      url: `http://localhost:8080/games/${game_id}/guesses`,
      method: 'POST',
      body: JSON.stringify(body),
      andThen: callback
    }
  })
}))(Mobile);
