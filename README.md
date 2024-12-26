## Developed by Romer Necesario

## Catalog Services Overview
This repository contains the Catalog Services for managing product data, inventory, and synchronization with an Ecommerce site. It is designed to handle catalog operations, including importing, updating, and exporting product information.

**Features:**

1. Product Catalog Management: Add, update, and delete product entries.
2. Inventory Synchronization: Sync product data with third-party platforms.
3. Automated Reports: Generate sales and inventory reports.
4. API Integration: Interface with external services via REST API.

**Edit a file, create a new file, and clone from Bitbucket in under 2 minutes**

When you're done, you can delete the content in this README and update the file with details for others getting started with your repository.

*We recommend that you open this README in another tab as you perform the tasks below. You can [watch our video](https://youtu.be/0ocf7u76WSo) for a full demo of all the steps in this tutorial. Open the video in a new tab to avoid leaving Bitbucket.*

---

## Edit a file

You’ll start by editing this README file to learn how to edit a file in Bitbucket.

1. Click **Source** on the left side.
2. Click the README.md link from the list of files.
3. Click the **Edit** button.
4. Delete the following text: *Delete this line to make a change to the README from Bitbucket.*
5. After making your change, click **Commit** and then **Commit** again in the dialog. The commit page will open and you’ll see the change you just made.
6. Go back to the **Source** page.

---

## Create a file

Next, you’ll add a new file to this repository.

1. Click the **New file** button at the top of the **Source** page.
2. Give the file a filename of **contributors.txt**.
3. Enter your name in the empty file space.
4. Click **Commit** and then **Commit** again in the dialog.
5. Go back to the **Source** page.

Before you move on, go ahead and explore the repository. You've already seen the **Source** page, but check out the **Commits**, **Branches**, and **Settings** pages.

---

## Clone a repository

Use these steps to clone from SourceTree, our client for using the repository command-line free. Cloning allows you to work on your files locally. If you don't yet have SourceTree, [download and install first](https://www.sourcetreeapp.com/). If you prefer to clone from the command line, see [Clone a repository](https://confluence.atlassian.com/x/4whODQ).

1. You’ll see the clone button under the **Source** heading. Click that button.
2. Now click **Check out in SourceTree**. You may need to create a SourceTree account or log in.
3. When you see the **Clone New** dialog in SourceTree, update the destination path and name if you’d like to and then click **Clone**.
4. Open the directory you just created to see your repository’s files.

Now that you're more familiar with your Bitbucket repository, go ahead and add a new file locally. You can [push your change back to Bitbucket with SourceTree](https://confluence.atlassian.com/x/iqyBMg), or you can [add, commit,](https://confluence.atlassian.com/x/8QhODQ) and [push from the command line](https://confluence.atlassian.com/x/NQ0zDQ).

## Setup a project

**Step 1: Install Go**
Ensure that Go is installed on your system. You can check if Go is installed by running:
**go version**

If Go is not installed, download and install it from the official Go website: [-> (https://golang.org/dl/)].

**Step 2: Initialize a Go module in this directory:**
**go mod init catalog-service**
This will create a go.mod file that tracks your dependencies.

**Step 3: Install Dependencies (if any)**
If your main.go file uses any external packages (such as the logger package in the example), make sure to install them:
**go get <package-name>**

**Step 4: Run the Application**
To run the main.go file directly, use the following command:
**go run cmd/<specific-subservices>/main.go**

**Step 5: Build the Application**
To build your Go application into an executable file, run the following command:
**go build -o cmd/<specific-subservices>**

This will create an executable file named catalog-service in your project directory. You can then run the application directly using:
**./cmd/<specific-subservices>**
