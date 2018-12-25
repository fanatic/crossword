import React, { Component } from 'react';
import { connect, PromiseState } from 'react-refetch';
import PropTypes from 'prop-types';
import Board from './Board';
import './TV.css';

class TV extends Component {
  static propTypes = {
    getGame: PropTypes.instanceOf(PromiseState).isRequired
  };

  render() {
    const { game } = this.props;

    return (
      <div>
        <h2>TV Play</h2>
        <div className="tv-layout">
          <div className="tv-left">
            <p>
              <className="score-header">
               Score
              <ul className="scoreboard">
                <li>Jason: 10000 pts
                </li>
                <li>Heather: 5000 pts</li>
              </ul>
            </p>
            <p>
              <className="last-clue-header">
                Last Clue:{' '}
              <br />
              <span className="last-clue">
              >
                (Clue 1 Across: Animal)
              </span>
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
            <span className="countdown"> }>:30 </span>
            <p className="game-id">           
              <button is="google-cast-button" />
              <strong>Game ID: BLAH</strong>
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

export default connect(props => ({}))(TV);
