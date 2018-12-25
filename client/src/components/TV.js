import React, { Component } from 'react';
import { connect, PromiseState } from 'react-refetch';
import PropTypes from 'prop-types';
import Board from './Board';

class TV extends Component {
  static propTypes = {
    getGame: PropTypes.instanceOf(PromiseState).isRequired
  };

  render() {
    const { game } = this.props;

    return (
      <div>
        <h2>TV Play</h2>
        <div style={{ display: 'flex', fontSize: '34px' }}>
          <div style={{ flex: '1' }}>
            <p>
              <strong>Score</strong>
              <ul
                style={{ textAlign: 'left', width: '200px', margin: 'auto', listStyleType: 'none', fontSize: '20px' }}
              >
                <li>Jason: 10000 pts</li>
                <li>Heather: 5000 pts</li>
              </ul>
            </p>
            <p>
              <strong>Last Clue (Clue 1 Across: Animal)</strong>
              <ul
                style={{ textAlign: 'left', width: '200px', margin: 'auto', listStyleType: 'none', fontSize: '20px' }}
              >
                <li>Jason: DOG (0 pts)</li>
                <li>Heather: CAT (10 pts)</li>
              </ul>
            </p>
            <p style={{ fontSize: '34px' }}>
              <strong>
                Clue 2 Across:
                <br />
                TV dinner guest
              </strong>
            </p>
            <span style={{ textAlign: 'left', border: '1px black solid', padding: '5px' }}>:30</span>
            <p>
              <button is="google-cast-button" />
              <strong>Game ID: BLAH</strong>
            </p>
          </div>
          <div style={{ flex: '2' }}>
            <Board game={game} />
          </div>
        </div>
      </div>
    );
  }
}

export default connect(props => ({}))(TV);
