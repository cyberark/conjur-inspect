#!/usr/bin/env groovy
@Library("product-pipelines-shared-library") _

// Automated release, promotion and dependencies
properties([
  // Include the automated release parameters for the build
  release.addParams(),
])

// Performs release promotion.  No other stages will be run
if (params.MODE == "PROMOTE") {
  release.promote(params.VERSION_TO_PROMOTE) { sourceVersion, targetVersion, assetDirectory ->
    // Any assets from sourceVersion Github release are available in assetDirectory
    // Any version number updates from sourceVersion to targetVersion occur here
    // Any publishing of targetVersion artifacts occur here
    // Anything added to assetDirectory will be attached to the Github Release
  }

  // Copy Github Enterprise release to Github
  release.copyEnterpriseRelease(params.VERSION_TO_PROMOTE)
  return
}

pipeline {
  agent { label 'conjur-enterprise-common-agent' }

  options {
    timestamps()
    buildDiscarder(logRotator(numToKeepStr: '30'))
    timeout(time: 2, unit: 'HOURS')
  }

  triggers {
    cron(getDailyCronString())
  }

  environment {
    // Sets the MODE to the specified or autocalculated value as appropriate
    MODE = release.canonicalizeMode()
  }

  stages {
    // Aborts any builds triggered by another project that wouldn't include any changes
    stage ("Skip build if triggering job didn't create a release") {
      when {
        expression {
          MODE == "SKIP"
        }
      }
      steps {
        script {
          currentBuild.result = 'ABORTED'
          error("Aborting build because this build was triggered from upstream, but no release was built")
        }
      }
    }

    stage('Scan for internal URLs') {
      steps {
        script {
          detectInternalUrls()
        }
      }
    }

    stage('Get InfraPool Agent') {
      steps {
        script {
          INFRAPOOL_EXECUTORV2_AGENT_0 = getInfraPoolAgent.connected(type: "ExecutorV2", quantity: 1, duration: 1)[0]
        }
      }
    }

    // Generates a VERSION file based on the current build number and latest version in CHANGELOG.md
    stage('Validate Changelog and set version') {
      steps {
        script {
          INFRAPOOL_EXECUTORV2_AGENT_0.agentSh './bin/parse-changelog'
          updateVersion(INFRAPOOL_EXECUTORV2_AGENT_0, "CHANGELOG.md", "${BUILD_NUMBER}")
        }
      }
    }

    stage('Get latest upstream dependencies') {
      steps {
        updateGoDependencies(INFRAPOOL_EXECUTORV2_AGENT_0, '${WORKSPACE}/go.mod')
      }
    }

    stage('Run Unit Tests') {
      steps {
        script {
          INFRAPOOL_EXECUTORV2_AGENT_0.agentSh './bin/test-unit'
          INFRAPOOL_EXECUTORV2_AGENT_0.agentStash name: 'junit-coverage', includes: '*.xml'
        }
      }

      post {
        always {
          unstash 'junit-coverage'
          junit 'junit.xml'

          cobertura autoUpdateHealth: false,
            autoUpdateStability: false,
            coberturaReportFile: 'coverage.xml',
            conditionalCoverageTargets: '70, 0, 0',
            failUnhealthy: false,
            failUnstable: false,
            maxNumberOfBuilds: 0,
            lineCoverageTargets: '70, 0, 0',
            methodCoverageTargets: '70, 0, 0',
            onlyStable: false,
            sourceEncoding: 'ASCII',
            zoomCoverageChart: false
            codacy action: 'reportCoverage', filePath: "coverage.xml"
        }
      }
    }

    // This produces the conjur-inspect binaries for integration tests and
    // pushing a release when this is a RELEASE build.
    stage('Create Release Assets') {
      steps {
        script {
          INFRAPOOL_EXECUTORV2_AGENT_0.agentSh "bin/build-release"
        }
      }
      post {
        always {
          script {
            INFRAPOOL_EXECUTORV2_AGENT_0.agentArchiveArtifacts artifacts: 'dist/*.tar.gz', allowEmptyArchive: true, fingerprint: false
            INFRAPOOL_EXECUTORV2_AGENT_0.agentArchiveArtifacts artifacts: 'dist/*.rpm', allowEmptyArchive: true, fingerprint: false
            INFRAPOOL_EXECUTORV2_AGENT_0.agentArchiveArtifacts artifacts: 'dist/*.deb', allowEmptyArchive: true, fingerprint: false
          }
        }
      }
    }

    // Currently the integration tests don't pass or fail the build based on
    // any conditions. Rather, that provide an easy way to inspect the result
    // from a few various environments without running the tool manually.
    stage('Run Integration Tests') {
      steps {
        script {
          INFRAPOOL_EXECUTORV2_AGENT_0.agentSh 'bin/test-integration'
        }
      }
      post {
        always {
          script {
            INFRAPOOL_EXECUTORV2_AGENT_0.agentArchiveArtifacts artifacts: 'ci/integration/results/**', allowEmptyArchive: true, fingerprint: false
          }
        }
      }
    }

    stage('Release') {
      when {
        expression {
          MODE == "RELEASE"
        }
      }

      steps {
        script {
          release(INFRAPOOL_EXECUTORV2_AGENT_0) { billOfMaterialsDirectory, assetDirectory, toolsDirectory ->
            // Publish release artifacts to all the appropriate locations
            // Copy any artifacts to assetDirectory to attach them to the Github release

            // Create Go application SBOM using the go.mod version for the golang container image
            INFRAPOOL_EXECUTORV2_AGENT_0.agentSh """export PATH="${toolsDirectory}/bin:${PATH}" && go-bom --tools "${toolsDirectory}" --go-mod ./go.mod --image "golang" --main "cmd/conjur-inspect/" --output "${billOfMaterialsDirectory}/go-app-bom.json" """
            // Create Go module SBOM
            INFRAPOOL_EXECUTORV2_AGENT_0.agentSh """export PATH="${toolsDirectory}/bin:${PATH}" && go-bom --tools "${toolsDirectory}" --go-mod ./go.mod --image "golang" --output "${billOfMaterialsDirectory}/go-mod-bom.json" """

            // Add goreleaser artifacts to release
            INFRAPOOL_EXECUTORV2_AGENT_0.agentSh """cp dist/*.tar.gz "${assetDirectory}" """
            INFRAPOOL_EXECUTORV2_AGENT_0.agentSh """cp dist/*.rpm "${assetDirectory}" """
            INFRAPOOL_EXECUTORV2_AGENT_0.agentSh """cp dist/*.deb "${assetDirectory}" """
            INFRAPOOL_EXECUTORV2_AGENT_0.agentSh """cp "dist/SHA256SUMS.txt" "${assetDirectory}" """
          }
        }
      }
    }
  }

  post {
    always {
      releaseInfraPoolAgent(".infrapool/release_agents")
    }
  }
}
