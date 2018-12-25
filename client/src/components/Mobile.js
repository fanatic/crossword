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

    this.props.postGuess(this.props.game.id, { player_name: this.props.player_name, guess: this.state.guess });
  };

  render() {
    const { player_name, game } = this.props;

    const all_players = (game.current_players || []).map(p => p.name);
    const guessed_players = ((game.current_clue || {}).guesses || []).map(g => g.player.name);
    const missing_players = all_players - guessed_players || [];

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
              <input
                type="text"
                name="guess"
                value={this.state.guess}
                onChange={this.handleInputChange}
                minLength={game.current_clue.answer.length}
                maxLength={game.current_clue.answer.length}
              />
              <input type="submit" value="Guess" />
            </form>
            <br />
            <small>Waiting on {missing_players.join(', ')}.</small>
          </React.Fragment>
        )}
      </div>
    );
  }
}

export default connect(props => ({
  postGuess: (game_id, body) => ({
    postGuessResponse: {
      url: `http://192.168.3.38:8080/games/${game_id}/guesses`,
      method: 'POST',
      body: JSON.stringify(body),
      force: true,
      refreshing: true
    }
  })
}))(Mobile);
