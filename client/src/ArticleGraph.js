import Article from './Article';

export default class ArticleGraph {
  constructor(rootId) {
    this.rootId = rootId;
    this.citations = [];
    this.root = null;
    this.idToNode = new Map();
    this.nodeToDepth = new Map();
  }

  // Methods for fetching from Semantic Scholar API.
  fetch = async () => {
    this.root = new Article(this.rootId);
    await this.root.fetch();
    this.idToNode.set(this.rootId, this.root);
    this.depthToNodes = [[this.root]];
  }

  fetchNode = async (id) => {
    if (this.idToNode.has(id)) {
      return this.idToNode.get(id);
    }
    const node = new Article(id);
    await node.fetch();
    this.idToNode.set(id, node);
    return node;
  };

  fetchCitingArticles = async (id) => {
    const article = this.idToNode.get(id);
    const citingArticles = await article.fetchCitingArticles();
    for (let article of citingArticles) {
      if (!this.idToNode.has(article.id)) {
        this.idToNode.set(article.id, article)
      }
    }
  };

  // TODO: this eventually needs to be paged
  fetchNextLayer = async () => {
    const maxDepth = this.depthToNodes.length;
    const lastLayer = this.depthToNodes[maxDepth - 1];
    const citingArticles = (await Promise.all(
      lastLayer
        .map(article => article.fetchCitingArticles())
    )).flat();
    const nextLayer = [];
    for (let article of citingArticles) {
      if (!this.idToNode.has(article.id)) {
        this.idToNode.set(article.id, article);
        nextLayer.push(article);
      }
    }
    this.depthToNodes.push(nextLayer);
    return this;
  }

  // Methods for exploring the graph.
  getArticlesForDepth = (depth) => {
  }

  getLayers = () => {
    return this.depthToNodes;
  }

  getDisjointSubtrees = () => {
  }

  getDepthOfChild = (id) => {
  }

}
