import React, { Component } from "react";
import ArticleGraph from "./ArticleGraph";
import "./App.css";

class App extends Component {
  state = {};

  async componentDidMount() {
    console.log("hit");
    const response = await fetch(
      "http://localhost:8000/article?paperId=3852968082a16db8be19b4cb04fb44820ae823d4&depth=2"
    );
    const data = await response.json();
    this.setState({ data });
  }

  render() {
    const { data } = this.state;
    const renderedLayers = [];
    if (data) {
      for (let layer of Object.values(data)) {
        const renderedLayer = [];
        for (let article of layer) {
          renderedLayer.push(<div key={article.id}>{article.title}</div>);
        }
        renderedLayers.push(<div>{renderedLayer}</div>);
      }
      return <div style={{ display: "flex" }}>{renderedLayers}</div>;
    }
    return <div>Loading</div>;
  }
}

export default App;
