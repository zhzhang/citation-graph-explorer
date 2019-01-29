import falcon
from resources.article import ArticleResource

api = falcon.API()
api.add_route('/article', ArticleResource())
