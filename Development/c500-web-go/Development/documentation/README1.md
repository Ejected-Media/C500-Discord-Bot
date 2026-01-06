This is the crucial step where everything connects. We have the backend logic (Go), the styling (CSS), and the assets (images), but we are missing the actual HTML structure that ties them together and displays them to the user.
Therefore, the next logical step is to create the HTML templates.
I will provide the source code for two essential templates:
 * templates/base.html: The master layout file that includes your CSS and defines the common structure (header, footer).
 * templates/docs/buyer.html: An example content page showing how to embed one of your diagrams using the CSS classes we just created.
1. The Master Layout (templates/base.html)
This file defines the outer shell of your website. Notice how it links the two CSS files in the <head> and defines a {{block "content" .}}{{end}} area where individual pages will inject their specific information.
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{block "title" .}}c500-web-go{{end}}</title>
    
    <link rel="stylesheet" href="/static/css/style.css">
    <link rel="stylesheet" href="/static/css/docs.css">
</head>
<body>
    <header>
        <nav>
            <div class="logo">c500-web-go</div>
            <ul>
                <li><a href="/">Home</a></li>
                <li><a href="/docs/">Documentation</a></li>
            </ul>
        </nav>
    </header>

    <main>
        {{block "content" .}}{{end}}
    </main>

    <footer>
        <p>&copy; 2026 c500-web-go Project. All rights reserved.</p>
    </footer>
</body>
</html>

2. The Buyer's Manual Page (templates/docs/buyer.html)
This template "extends" the base layout. It fills in the title block and the content block.
Pay close attention to the <figure> element. This is exactly how you utilize the CSS we just wrote (figure.diagram-figure) to make your diagrams look professional and responsive. You would repeat this pattern for the seller, developer, and admin pages, just changing the content and the image source.
{{extends "base.html"}}

{{define "title"}}Buyer's Manual - c500-web-go Docs{{end}}

{{define "content"}}
<article class="docs-content">
    <h1>Buyer's Documentation Manual</h1>
    <p class="intro">Welcome to the buyer's guide. This page outlines the simple, secure process of purchasing items through our Discord integration and web checkout system.</p>

    <h2>The Purchase Workflow</h2>
    <p>The buying process is designed to be seamless, taking you directly from a Discord notification to a secure payment page.</p>

    <figure class="diagram-figure">
        <img src="/static/images/buyer-flow.png" alt="Diagram illustrating the workflow from a Discord 'Buy Now' link to Stripe payment confirmation.">
        <figcaption>Figure 1: Visual overview of the buyer purchase journey.</figcaption>
    </figure>

    <h3>Step-by-Step Guide</h3>
    <ol>
        <li><strong>Locate an Item:</strong> Find an item for sale in the designated Discord channel.</li>
        <li><strong>Click "Buy Now":</strong> Click the link provided in the bot's message embed.</li>
        <li><strong>Secure Checkout:</strong> You will be redirected to a secure c500-web-go checkout page in your browser. Review order details.</li>
        <li><strong>Complete Payment:</strong> Click "Pay with Stripe" to finalize the transaction securely.</li>
        <li><strong>Confirmation:</strong> You will receive an on-screen confirmation and a receipt via email.</li>
    </ol>
</article>
{{end}}

After this, the final remaining major piece is the main.go file to tie everything together and run the server.
