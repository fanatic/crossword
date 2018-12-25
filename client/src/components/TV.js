import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Board from './Board';
import './TV.css';

export default class TV extends Component {
  static propTypes = {
    game: PropTypes.object.isRequired
  };

  render() {
    const { game } = this.props;

    return (
      <div>
        <h2>TV Play</h2>
        <div className="tv-layout">
          <div className="tv-left">
            <p>
              <span className="score-header">Score</span>
              <ul className="scoreboard">
                <li>Jason: 10000 pts</li>
                <li>Heather: 5000 pts</li>
              </ul>
            </p>
            <p>
              <span className="last-clue-header">
                Last Clue: <br />
              </span>
              <span className="last-clue">> (Clue 1 Across: Animal)</span>
              <ul className="last-clue-score">
                <li>Jason: DOG (0 pts)</li>
                <li>Heather: CAT (10 pts)</li>
              </ul>
            </p>
            <p className="current-clue">
              Clue 2 Across:
              <br />
              TV dinner guest
            </p>
            <span className="countdown"> :30 </span>
            <p className="game-id">
              <button is="google-cast-button" />
              Game ID: BLAH
            </p>
          </div>
          <div className="board-layout">
            <Board game={game} />
          </div>
        </div>
      </div>
    );
  }
}
