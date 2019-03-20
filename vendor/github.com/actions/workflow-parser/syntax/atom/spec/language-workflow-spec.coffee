describe "Workflow grammar", ->
  grammar = null

  beforeEach ->
    waitsForPromise ->
      atom.packages.activatePackage("language-workflow")

    runs ->
      grammar = atom.syntax.grammarForScopeName("source.workflow")

  it "parses the grammar", ->
    expect(grammar).toBeTruthy()
    expect(grammar.scopeName).toBe "source.workflow"
