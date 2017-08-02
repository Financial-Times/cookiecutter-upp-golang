import os


def remove_file(file_name):
    if os.path.exists(file_name):
        os.remove(file_name)

# Get the root project directory
PROJECT_DIRECTORY = os.path.realpath(os.path.curdir)

sample_endpoint = '{{ cookiecutter.add_sample_http_endpoint }}'

if sample_endpoint == "no":
    file_name = os.path.join(PROJECT_DIRECTORY, "handler.go")
    remove_file(file_name)


