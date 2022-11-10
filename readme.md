# PANDA Studio

Prerequisites:
- For building/running the web ui, you will need node 14.15.1.
- You may also need node 14.13.1 to install dependencies. If 14.15.1 doesn't work, switch to 14.13.1 to install dependencies then
switch back to 14.15.1 to build.
- With node 14.13.1, run "npm install --force" to install necessary dependencies
- Python 3.10 for the server project.
- Python 3.7 for the backend project.
- You should create a separate venv in each of these projects to install respective dependencies. A tutorial on how to do
so can be found here: https://docs.python.org/3/library/venv.html
- Once inside the active venv, you will need to pip install flask and flask_cors for the server project
- You will also need to pip install pandare inside the backend project

Running the web ui:
- Install the necessary dependencies (described above)
- Ensure that you are in the panda-studio-ui directory
- Run "npm run dev" to start the web client

Running the web server:
- Make sure necessary dependencies are installed
- Run "flask --app server.py run"
- The web server needs to be running in order for the UI to properly call to the API

Running the backend process:
- Don't. It just exists. Something will call this code to run it, but the dependencies do need to exist.



