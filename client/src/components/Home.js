import React, { Component } from 'react';
import './Home.css';

class Home extends Component {
  render() {
    return (
      <div className="App">
        <h2>Home Screen</h2>
        <p>
          <input type="submit" value="Host Game" />
        </p>
        <p>
          <input type="text" placeholder="Name" />
          <input type="text" placeholder="Game ID" />
          <input type="submit" value="Join Game" />
        </p>
        <hr />
        <h2>Mobile Play</h2>
        <p>Clue 1 Across: TV dinner guest</p>
        <p className="clueBox">
          <table>
            <tr>
              <td style={{ border: '1px black solid' }}>&nbsp;</td>
              <td style={{ border: '1px black solid' }}>A</td>
              <td style={{ border: '1px black solid' }}>&nbsp;</td>
            </tr>
          </table>
        </p>
        <input type="text" />
        <input type="submit" value="Guess" />
        <br />
        <small>You guessed CAT. Waiting on Heather, Melanie.</small>
        <hr />
        <h2>TV Play</h2>
        <p>
          <strong>Score</strong>
          <ul style={{ textAlign: 'left', width: '200px', margin: 'auto', listStyleType: 'none' }}>
            <li>Jason: 10000 pts</li>
            <li>Heather: 5000 pts</li>
          </ul>
        </p>
        <p>
          <strong>Last Clue (Clue 1 Down: Animal)</strong>
          <ul style={{ textAlign: 'left', width: '200px', margin: 'auto', listStyleType: 'none' }}>
            <li>Jason: DOG (0 pts)</li>
            <li>Heather: CAT (10 pts)</li>
          </ul>
        </p>
        <p>Board here</p>
        <p>Clue 1 Across: TV dinner guest</p>
        <span style={{ textAlign: 'left', border: '1px black solid' }}>:30</span>
      </div>
    );
  }
}

export default Home;
