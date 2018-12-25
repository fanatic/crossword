import React, { Component } from 'react';
import './Board.scss';

class BoardFragment extends Component {
  render() {
    const { clue } = this.props;

    let rows = 1;
    let cols = 1;
    if (clue.direction === 'across') {
      cols = clue.answer.length;
    } else {
      rows = clue.answer.length;
    }
    let items = [];
    for (let i = 0; i < clue.answer.length; i++) {
      let value = '';
      if (clue.answer[i] !== '?') {
        value = clue.answer[i];
      }
      items.push(
        <input
          key={i}
          className="crossword-board__item"
          type="text"
          minLength="1"
          maxLength="1"
          required="required"
          value={value}
          readOnly
        />
      );
    }

    const crosswordBoardStyle = {
      width: 50 * cols + 'px',
      height: 50 * rows + 'px',
      gridTemplateRows: 'repeat(' + rows + ', ' + 100 / rows + '%)',
      gridTemplateColumns: 'repeat(' + cols + ', ' + 100 / cols + '%)'
    };

    return (
      <div>
        <div className="crossword-board-container">
          <div className="crossword-board" style={crosswordBoardStyle}>
            {items}

            <div className="crossword-board crossword-board--labels" style={crosswordBoardStyle}>
              <span
                style={{
                  gridColumnStart: 1,
                  gridColumnEnd: 1,
                  gridRowStart: 1,
                  gridRowEnd: 1
                }}
                className={`crossword-board__item-label`}
              >
                <span className="crossword-board__item-label-text">{clue.number}</span>
              </span>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default BoardFragment;
