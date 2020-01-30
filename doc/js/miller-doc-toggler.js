// TO DO:
// * module.exports
// * apply to mlrdoc (copypasta & git-ci probably ...)
// * document this class

'use strict';

class MillerDocToggler {
  // ----------------------------------------------------------------
  // PUBLIC METHODS

  // Prefix for toggleable div names, without leading '#'
  constructor(toggleableDivPrefix, buttonSelectFontWeight, buttonDeselectFontWeight) {
    this._allDivNames = [];
    const divs = document.querySelectorAll('div');
    for (let div of divs) {
      const id = div.id;
      if (id.startsWith(toggleableDivPrefix)) {
        this._allDivNames.push(id);
      }
    }

    this._buttonSelectFontWeight = buttonSelectFontWeight;
    this._buttonDeselectFontWeight = buttonDeselectFontWeight;
    this._allExpanded = false;
  }

  // Opening one closes others, unless expand-all.
  //
  // * If everything is expanded, selecting one means *keep* it expanded
  //   but collapse everything else.
  //
  // * If only one is expanded, then:
  //   o selecting that same one means collapse it;
  //   o selecting that another means collapse the old one and expand
  //     the new one.
  expandUniquely = (divName) => {
    const eleDiv = document.getElementById(divName);
    const button = document.getElementById(divName+"_button")
    if (eleDiv != null) {
      if (this._allExpanded) {
        this.collapseAll();
        if (button != null) {
          this._makeButtonSelected(button);
        }
        eleDiv.style.display = "block";
      } else {
        const state = eleDiv.style.display;
        this.collapseAll();
        if (state === "block") {
          this._makeButtonDeselected(button);
          eleDiv.style.display = "none";
        } else {
          if (button != null) {
            this._makeButtonSelected(button);
          }
          eleDiv.style.display = "block";
        }
      }
    }
    this._allExpanded = false;
  };

  expandAll = () => {
    for (let divName of this._allDivNames) {
      this._expand(divName);
    }
    this._allExpanded = true;
  };

  collapseAll = () => {
    for (let divName of this._allDivNames) {
      this._collapse(divName);
    }
    this._allExpanded = false;
  }

  toggle = (divName) => {
    const div = document.getElementById(divName);
    if (div != null) {
      const state = div.style.display;
      if (state == 'block') {
        div.style.display = 'none';
      } else {
        div.style.display = 'block';
      }
    }
  }

  // ----------------------------------------------------------------
  // PRIVATE METHODS

  _expand = (divName) => {
    const eleDiv = document.getElementById(divName);
    const button = document.getElementById(divName+"_button")
    if (eleDiv != null) {
      eleDiv.style.display = "block";
    }
    if (button != null) {
      this._makeButtonSelected(button)
    }
  };

  _collapse = (divName) => {
    const eleDiv = document.getElementById(divName);
    const button = document.getElementById(divName+"_button")
    if (eleDiv != null) {
      eleDiv.style.display = "none";
    }
    if (button != null) {
      this._makeButtonDeselected(button)
    }
  };

  _makeButtonSelected = (button) => {
    button.style.fontWeight = 'bold';
    //button.style.borderWidth = 'thin';
    //button.style.borderStyle = 'solid';
  };

  _makeButtonDeselected = (button) => {
    button.style.fontWeight = 'normal';
    //button.style.borderWidth = 'none';
    //button.style.borderStyle = 'none';
  };

}

// module.exports = MillerDocToggler;
