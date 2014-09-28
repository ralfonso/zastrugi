/** @jsx React.DOM */

var app = app || {};

(function(){
  'use strict';

  var Zastrugi = React.createClass({
    render: function() {
      return <div>hello!</div>;
    }
  });

  function init() {
    React.renderComponent(<Zastrugi/>, document.getElementById('container'));
  }

  init();
})();
