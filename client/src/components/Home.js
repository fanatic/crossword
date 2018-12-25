import React, { Component } from 'react';
import { connect, PromiseState } from 'react-refetch';
import PropTypes from 'prop-types';

class Home extends Component {
  static propTypes = {
    updateGame: PropTypes.func.isRequired,
    updateGameResponse: PropTypes.instanceOf(PromiseState),
    postGame: PropTypes.func.isRequired,
    postGameResponse: PropTypes.instanceOf(PromiseState)
  };

  render() {
    return (
      <div>
        <h2>Home Screen</h2>
        <p>
          <input type="submit" value="Host Game" />
        </p>
        <p>
          <input type="text" placeholder="Name" />
          <input type="text" placeholder="Game ID" />
          <input type="submit" value="Join Game" />
        </p>
      </div>
    );
  }
}

export default connect(props => ({
  updateGame: body => ({
    updateGameResponse: {
      url: `/games/${props.gameID}`,
      method: 'PUT',
      body: JSON.stringify(body)
    }
  }),
  postGame: body => ({
    postGameResponse: {
      url: `/games`,
      method: 'POST',
      body: JSON.stringify(body)
    }
  })
}))(Home);
