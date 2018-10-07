Feature: Building an App
  I want to build a new release of an application

  Background:
    Given there is an app named toast
    And there is a shared configuration for the toast app
    And I have a cached copy of the toast app
    And there is a buildpack installed that can build toast
    And there is a buildpack installed that cannot build toast

  Scenario: Default Behavior
    When I run `ogun build toast`
    Then the shared config is applied to the build
    And the proper buildpack is detected
    And a new toast slug is generated



