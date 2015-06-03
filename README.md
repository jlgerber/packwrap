# Packwrap - Package/Wrapper system

I am going to say up front that this is a simple DCC wrapper system circa 2000. It would have been written in Perl, and it would have been more full featured than this. However, it is a fairly mindless vehicle to play with Golang. And it is fine for home...

## Introduction

Packwrap is designed to solve a couple of common problems in VFX. Namely, it aims to handle environment initialization, application invocation, tool chain creation.

### Environment Initialization

The first problem packwrap tackles is coming up with a system for initializing the application environment. Since most DCCs can and need to be customized via scores of environment variables, creating a tool for setting them is one of the first orders of business for any new pipeline. Historically, this has been accomplished via in a scripting language - either sh,  Python or Perl if you go back far enough. 

## Components

* executable
* manifest

### example execution

paw run maya 2014.0.0
paw run  houdini 14.0.335
paw list
paw versions maya

### example structure

manifest/
    modo/
        801.0.0/manifest.yaml
	    901.0.0/manifest.yaml
	houdiini/
	    14.0.335/manifest.yaml