# Packwrap - Package/Wrapper system

## Introduction

Packwrap is designed to solve a couple of common problems in VFX. Namely, it aims to handle environment initialization, application invocation, tool chain creation.

### Environment Initialization

The first problem packwrap tackles is coming up with a system for initializing the application environment. Since most DCCs can and need to be customized via scores of environment variables, creating a tool for setting them is one of the first orders of business for any new pipeline. Historically, this has been accomplished via in a scripting language - either sh,  Python or Perl if you go back far enough. 

## Components

* executable
* manifest

### example execution

packwrap run maya
pacwrap run mplayer --context houdini-14
packwrap list
packwrap versions maya

### example structure

manifest/
    modo/
        801.0.0/manifest.yaml
	    901.0.0/manifest.yaml
	houdiini/
	    14.0.335/manifest.yaml