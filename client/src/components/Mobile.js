import React, { Component } from 'react';
import { connect, PromiseState } from 'react-refetch';
import PropTypes from 'prop-types';
import BoardFragment from './BoardFragment';

class Mobile extends Component {
  static propTypes = {
    getGame: PropTypes.instanceOf(PromiseState).isRequired
  };

  render() {
    const { game } = this.props;

    return (
      <div>
        <h2>Mobile Play</h2>
        <p>{game.current_clue.description}</p>
        <BoardFragment clue={game.current_clue} />
        <input type="text" />
        <input type="submit" value="Guess" />
        <br />
        <small>Waiting on Heather, Melanie.</small>
      </div>
    );
  }
}

export default connect(props => ({
  getGame: ``
}))(Mobile);
