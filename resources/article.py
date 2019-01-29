class ArticleResource:
    def on_post(self, req, resp):
        """Handles POST requests"""
        quote = {
            'quote': (
                "I've always been more interested in "
                "the future than in the past."
            ),
            'author': 'Grace Hopper'
        }
        print(req)

        resp.media = quote

    def on_get(self, req, resp):
        resp.body = "testing"
