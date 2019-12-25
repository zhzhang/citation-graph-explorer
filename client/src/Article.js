export default class Article {
  constructor(id) {
    this.id = id;
    this.citations = [];
  }

  fetch = async () => {
    const response = await fetch(
      `http://api.semanticscholar.org/v1/paper/${this.id}`
    );
    const json = await response.json();
    this.title = json.title;
    this.citationIds = [];
    for (let citation of json.citations) {
      this.citationIds.push(citation.paperId);
    }
    return this;
  };

  fetchCitingArticles = async () => {
    return await Promise.all(
      this.citationIds
        .map(id => new Article(id))
        .map(article => article.fetch())
    );
  };
}
