# Overview

The Actions Workflow language describes workflows, which map repository
events to sequences of actions that run in response to those events.  It
is based on Hashicorp HCL but supports only a subset of HCL features.
This document specifies that subset.

# By example

```
# Workflow files can have comments, which begin with a # character and
# continue to the end of the line.  Comments may be on their own line, or
# may appear to the right of real content.
#
# Workflow files can have a version specifier, which must appear before
# any other (non-blank, non-comment) content.  Currently, the only legal
# version is 0.
version = 0

# Workflow files contain one or more workflows, which map an event to one
# or more actions that the workflow resolves.  Each workflow has a name,
# which is a double-quoted string.  UTF-8 characters and C-style escapes
# are supported in the string.  The workflow name is displayed to the user
# but has no other form or meaning.  Following the workflow name is a
# block of key-value pairs, contained in curly braces.
workflow "this happens when I push" {
  # Each key-value pair is an identifier, an equals sign, and a value of
  # the correct type for that identifier.  Only specific identifiers are
  # allowed, and no key may appear twice in the same block.  For
  # workflows, those allowed identifiers are: on and resolves.  The "on"
  # key is required; "resolves" is optional.

  # "on" identifies the event that will cause Actions to run this
  # workflow.  It's value is a double-quoted string, case-insensitive,
  # drawn from the list of known event types.
  on = "fork"

  # "resolves" identifies one or more actions that will be resolved when
  # the given event occurs.  Resolving an action means running all
  # (transitive) prerequisites for that action and then running the action
  # itself.  Any failed action prevents subsequent actions from running.
  # The value for resolves may be either a string or an array of strings.
  # Arrays are designated with square brackets and commas.
  resolves = [ "goal1", "goal2" ]
}

# Workflow files also contain one or more actions.  Like workflows,
# actions have a double-quoted string name and a block of key-value pairs.
# The name for an action can be any user-selected string, which must match
# actions that workflows resolve.
action "goal1" {
  # The valid keys in an action block are: uses, needs, runs, args, env,
  # and secrets.  The uses key is required; all others are optional.

  # The "uses" keyword identifies what actual code this action will run.
  # The value is always a string, and may take three forms:
  #   - "./path", which identifies a directory in the current repository
  #     containing a Dockerfile that describes an action
  #   - "owner/repo/path@ref", which identifies a directory in another
  #     GitHub repository containing a Dockerfile that describes an
  #     action.  Ref may be a branch, tag, or SHA.
  #   - "docker://image", which identifies a docker image
  uses = "docker://alpine"
  # or: uses = "./local-directory"
  # or: uses = "actions/bin/filter@master"

  # The "needs" keyword identifies one or more actions that must complete
  # successfully before this action can begin.  The value can be a string
  # or an array of strings.
  needs = "ci"

  # The "runs" keyword identifies a command to run in the action's Docker
  # container.  Its value can be either a string or an array of strings.
  # If the value is a string, Actions will parse it by separating at
  # whitespace.  If the value is an array, the array elements are passed
  # literally to Docker with no further parsing.
  runs = "echo hello"

  # The "args" keyword identifies arguments to attach to whatever command
  # is running in the Docker container.  If the container has an
  # entrypoint, the arguments will be appended to it.  In the rare case
  # where an action specifies with "runs" and "args" keywords, their
  # contents are appended (runs, then args) and passed as the command to
  # run in the container.  The "args" value may be either a string or an
  # array.  If it is a string, it is parsed by separating at whitespace.
  args = [ "world" ]

  # The "env" keyword identifies environment variables that will be
  # present in the Docker container for this action.  The value is a hash
  # (object), surrounded by braces, mapping environment variable names to
  # environment variable values.  Both names and values are strings.
  # Environment variables beginning with "GITHUB_" are reserved.
  env = {
    KEY = "VALUE"
    KEY2 = "VALUE2"
  }

  # The "secrets" keyword identifies secrets that will be present as
  # environment variables in the container.  The values of these secrets
  # are stored elsewhere, not in .workflow files, so only the names of
  # secrets appear here.  The value is an array of strings.  No secret can
  # have the same name as an environment variable.  The number of secrets
  # allowed in an entire workflow file is currently limited to 100.
  secrets = [ "GITHUB_TOKEN" ]
}

# Each action named in a "resolves" or "needs" key must be present in the
# file.  Circular dependencies are prohibited.
action "ci" {
  uses = "./ci"
}

action "goal2" {
  uses = "docker://alpine"
  runs = "echo howdy"
}
```

# Grammar

The below is an [ANTLR4](https://github.com/antlr/antlr4) grammar specifying the Actions Workflow language. As a spec, it is likely not the best basis for a real parser. For instance, no effort has been made to make the grammar output intuitive errors.

```g4
grammar workflow;

workflow_file : version? (workflow | action)* ;

version : 'version' '=' INTEGER;

workflow : 'workflow' str '{' (on_kvp | resolves_kvp)* '}' ;

on_kvp : 'on' '=' event_string ;

resolves_kvp : 'resolves' '=' string_or_array ;

string_or_array : str | string_array ;

string_array : '[' (( str ',' )* str ','?)? ']' ;

action : 'action' str '{' action_kvps '}' ;

action_kvps : (uses_kvp | needs_kvp | runs_kvp | args_kvp | env_kvp | secrets_kvp)*;

uses_kvp : 'uses' '=' (DOCKER_USES | LOCAL_USES | REMOTE_USES) ;

needs_kvp : 'needs' '=' string_or_array ;

runs_kvp : 'runs' '=' string_or_array ;

args_kvp : 'args' '=' string_or_array ;

env_kvp : 'env' '=' '{' env_var* '}' ;

secrets_kvp : 'secrets' '=' ident_array ;

env_var : IDENTIFIER '=' str ;

ident_array : '[' ((QUOTED_IDENTIFIER ',')* QUOTED_IDENTIFIER ','?)? ']';

event_string : QUOTED_IDENTIFIER ;

str : QUOTED_IDENTIFIER | STRING;

// https://github.com/docker/distribution/blob/b75069ef13a1de846c0cdf964f5917f5b00c1a47/reference/reference.go
DOCKER_USES: '"docker://' (DOCKER_REGISTRY '/')? DOCKER_PATH_COMPONENT ('/' DOCKER_PATH_COMPONENT)* ( DOCKER_TAG |  DOCKER_DIGEST )? '"';

DOCKER_REGISTRY : HOST_COMPONENT ('.' HOST_COMPONENT)* (':' INTEGER)? ;

fragment DOCKER_PATH_COMPONENT : ALPHANUM+ ([._-] ALPHANUM+)*;
fragment HOST_COMPONENT : ALPHANUM | ALPHANUM [a-zA-Z0-9-]* ALPHANUM;

DOCKER_TAG : ':' [a-zA-Z0-9_]+ ;

DOCKER_DIGEST                            : '@' DIGEST_ALGORITHM ':' HEX+ ;
fragment DIGEST_ALGORITHM                : DIGEST_ALGORITHM_COMPONENT ( DIGEST_ALGORITHM_SEPERATOR DIGEST_ALGORITHM_COMPONENT )*;
fragment DIGEST_ALGORITHM_SEPERATOR      : [+.-_];
fragment DIGEST_ALGORITHM_COMPONENT      : [A-Za-z] ALPHANUM*;

LOCAL_USES : '"./' SAFECODEPOINT* '"';

REMOTE_USES : '"' GITHUB_OWNER '/' GITHUB_REPO ('/' ~[/"]+)* '/'? '@' ~'/'? ~(["?*[ ^~:\\] | '\u0000'..'\u001F')+ ~([./])? '"';

// alphanums, and hyphens not at start or end
fragment GITHUB_OWNER : ALPHANUM+ ([a-zA-Z0-9\-]* ALPHANUM+)*?;
fragment GITHUB_REPO : [a-zA-Z0-9\-_.]+ ;

// before STRING to win on priority (all identifiers are valid strings)
QUOTED_IDENTIFIER : '"' IDENTIFIER '"';

IDENTIFIER : [a-zA-Z_] [a-zA-Z0-9_]*;

STRING
   : '"' ( ESC | SAFECODEPOINT )* '"'
   ;

fragment ESC
   : '\\' ( ["\\/bfnrt] )
   ;

fragment SAFECODEPOINT
   : ~ ["\\\u0000-\u001F\u007F]
   ;

LINE_COMMENT
    :   ('#' | '//') ~[\r\n]*
        -> skip
    ;

fragment ALPHANUM : [a-zA-Z0-9];
fragment HEX : [0-9a-fA-F]+;
INTEGER : [0-9]+;

WS : [\n \t\r] -> skip;
```
