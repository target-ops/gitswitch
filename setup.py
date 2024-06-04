from setuptools import setup, find_packages
setup(
    name='GitSwitch',
    version='0.0.1',
    packages=find_packages(include=['src', 'src.*'])
    entry_points={
        'console_scripts': ['GitSwitch=GitSwitch.src.main.py:main']
    }
)