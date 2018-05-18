#Python Installation Gotchas for Windows users
Download Python executable for Windows from [here](https://www.python.org/downloads/)  

If you run it with default configuration it installs in
C:\Users\{user.name}\AppData\Local\Programs\Python\Python36-32\python.exe  

In order to use *pip* you need to add to Environment Variables > Path  
`C:\Users\{user.name}\AppData\Local\Programs\Python\Python36-32\Scripts`  

When you run 
`pip install --user cookiecutter`

it will by default be installed in  
`C:\Users\{user.name}\AppData\Roaming\Python\Python36\Scripts`  

You will need to add this path to Environment Variables > Path as well.   