from setuptools import setup, find_packages

setup(
    name="GitSwitch",
    version="0.0.1",
    author="OfirHaim",
    author_email="ofir474@gmail.com",
    description="Manage multiple Git users for different vendors.",
    long_description=open('README.md').read(),
    long_description_content_type='text/markdown',
    url="https://github.com/target-ops/GitSwitch.git",
    packages=find_packages('src'),
    package_dir={'': 'src'},
    entry_points={
        'console_scripts': [
            'git-user-manager=main:main',
        ],
    },
    install_requires=[
        # List your dependencies here
    ],
    classifiers=[
        'Programming Language :: Python :: 3',
        'License :: OSI Approved :: MIT License',
        'Operating System :: OS Independent',
    ],
    python_requires='>=3.6',
)
