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
    And it exits successfully

  Scenario: With a custom buildpack
    Given the toast app has a custom buildpack
    When I run `ogun build toast`
    Then the custom buildpack is used to build the release
    And it exits successfully

  Scenario Outline: Providing a release name
    When I run `ogun build toast <Release Flag> 0987654321`
    Then a new toast slug named 0987654321 is generated
    And it exits successfully

    Examples:
      | Release Flag  |
      | -r            |
      | --release     |

    @failure
  Scenario: When buildpack detection fails
    Given there are no buildpacks that can build the app
    When I run `ogun build toast`
    Then I see an error regarding the lack of a viable buildpack
    And it exits with an error

    @failure
  Scenario: When app compilation fails
    Given there is an issue building the application
    When I run `ogun build toast`
    Then I see an error regarding the build failure
    And it exits with an error
