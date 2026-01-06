This is an exciting stage. You have successfully built the engine (the Go server), the chassis (the file structure and base template), and the paint job (CSS). You've even taken it for a test drive with the Buyer's Manual.
The immediate next step is to finish building out the remaining content pages.
Right now, if you try to visit /docs/seller, your server will panic because the templates/docs/seller.html file doesn't exist yet. We need to create the remaining three manual pages and the main documentation landing page.
Here is the plan to complete the documentation portal's foundation:
 * Create the Seller's Manual Template.
 * Create the Developer's Manual Template.
 * Create the Admin/Mod's Manual Template.
 * Create the Main Documentation Hub (Index) Template.
1. The Seller's Manual (templates/docs/seller.html)
Create this file. It follows the exact same pattern as the buyer template, but points to the seller diagram and uses seller-specific text.
Prerequisite: Ensure you have saved the second diagram (the seller flow) as static/images/seller-flow.png.
{{extends "base.html"}}

{{define "title"}}Seller's Manual - c500-web-go Docs{{end}}

{{define "content"}}
<article class="docs-content">
    <h1>Seller's Documentation Manual</h1>
    <p class="intro">This guide explains how to list items, manage your inventory, and track sales using the c500-web-go Discord bot and web dashboard.</p>

    <h2>Listing & Management Workflow</h2>
    <p>Managing your listings is done primarily through simple Discord commands, which interact with our central database.</p>

    <figure class="diagram-figure">
        <img src="/static/images/seller-flow.png" alt="Diagram illustrating the seller's workflow using Discord commands to manage listings.">
        <figcaption>Figure 2: The listing management process for sellers.</figcaption>
    </figure>

    <h3>Key Commands</h3>
    <ul>
        <li><strong>/list create:</strong> Start the process of adding a new item for sale. You will be prompted for details and images.</li>
        <li><strong>/manage-listings:</strong> View a list of your active items. From here you can edit prices or mark items as out of stock.</li>
        <li><strong>/sales-history:</strong> Get a DM report of your recent completed transactions.</li>
    </ul>
</article>
{{end}}

2. The Developer's Manual (templates/docs/developer.html)
Create this file. This page is more technical, focusing on the architecture diagram.
Prerequisite: Ensure the third diagram is saved as static/images/developer-arch.png.
{{extends "base.html"}}

{{define "title"}}Developer Docs - c500-web-go{{end}}

{{define "content"}}
<article class="docs-content">
    <h1>Developer's Documentation</h1>
    <p class="intro">A technical overview of the c500-web-go architecture, intended for contributors and integrators.</p>

    <h2>System Architecture</h2>
    <p>The system is built around a central Go API that manages data consistency between the Discord bot interface and the web payment integration.</p>

    <figure class="diagram-figure">
        <img src="/static/images/developer-arch.png" alt="Technical diagram showing the connection between the Go API, PostgreSQL database, Python Discord Bot, and Stripe API.">
        <figcaption>Figure 3: High-level system architecture and component interaction.</figcaption>
    </figure>

    <h3>Core Components</h3>
    <ul>
        <li><strong>c500-web-go API (Go):</strong> The RESTful backend service. It handles business logic, database transactions, and serves the frontend application.</li>
        <li><strong>PostgreSQL DB:</strong> The source of truth for users, listings, and order data.</li>
        <li><strong>Discord Bot (Python):</strong> Handles user interactions within Discord, forwarding commands to the Go API.</li>
        <li><strong>Stripe API:</strong> External payment processing used for secure checkout sessions.</li>
    </ul>
</article>
{{end}}

3. The Admin/Mod Manual (templates/docs/admin.html)
Create this file for your moderation team.
Prerequisite: Ensure the fourth diagram is saved as static/images/admin-dashboard.png.
{{extends "base.html"}}

{{define "title"}}Admin & Mod Manual - c500-web-go{{end}}

{{define "content"}}
<article class="docs-content">
    <h1>Admin & Moderator Manual</h1>
    <p class="intro">Tools and workflows for administrators and moderators to manage content, resolve disputes, and ensure platform safety.</p>

    <h2>The Admin Dashboard</h2>
    <p>Moderation is primarily handled through a secured web dashboard that provides a high-level view of platform activity.</p>

    <figure class="diagram-figure">
        <img src="/static/images/admin-dashboard.png" alt="Diagram showing the Admin User accessing the Web Dashboard to interact with the API and Database.">
        <figcaption>Figure 4: The moderator's workflow using the Admin Dashboard.</figcaption>
    </figure>

    <h3>Moderation Capabilities</h3>
    <ul>
        <li><strong>Flagged Listings Review:</strong> Review items flagged by users or automated systems for policy violations.</li>
        <li><strong>User Reports:</strong> Investigate reports filed against buyers or sellers.</li>
        <li><strong>Transaction Oversight:</strong> View details of specific transactions to assist with support inquiries.</li>
    </ul>
</article>
{{end}}

4. The Main Documentation Hub (templates/docs/index.html)
Finally, create the landing page for /docs/. This page needs to link to the four manuals you just created.
We already wrote a specific CSS class for this in the previous step: .docs-nav-list. Let's use it to create a nice grid of links.
{{extends "base.html"}}

{{define "title"}}Documentation Hub - c500-web-go{{end}}

{{define "content"}}
<section class="docs-hub">
    <h1>c500-web-go Documentation</h1>
    <p class="intro">Welcome to the official documentation hub. Select a guide below based on your role to get started.</p>

    <ul class="docs-nav-list">
        <li>
            <a href="/docs/buyer">
                <h3>Buyer's Manual</h3>
                <p>Learn how to purchase items securely.</p>
            </a>
        </li>
        <li>
            <a href="/docs/seller">
                <h3>Seller's Manual</h3>
                <p>Manage listings and view sales history.</p>
            </a>
        </li>
        <li>
            <a href="/docs/developer">
                <h3>Developer Docs</h3>
                <p>Technical architecture and API overview.</p>
            </a>
        </li>
        <li>
            <a href="/docs/admin">
                <h3>Admin & Mod Guide</h3>
                <p>Tools for moderation and platform management.</p>
            </a>
        </li>
    </ul>
</section>
{{end}}

Final Step
Once you have created these four files in your templates/docs/ directory, restart your Go server (go run main.go).
You can now navigate to http://localhost:8080/docs/ and you will see your hub page, with working links to all four fully illustrated documentation manuals.
