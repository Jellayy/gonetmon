/**
 * @type {import('semantic-release').GlobalConfig}
 */
module.exports = {
    branches: ["main"],
    plugins: [
        [
            '@semantic-release/commit-analyzer',
            {
                "preset": "angular",
                "parserOpts": {
                    "noteKeywords": ["BREAKING CHANGE", "BREAKING CHANGES", "BREAKING"]
                }
            }
        ],
        [
            '@semantic-release/release-notes-generator',
            {
                "preset": "angular",
                "parserOpts": {
                  "noteKeywords": ["BREAKING CHANGE", "BREAKING CHANGES", "BREAKING"]
                },
                "writerOpts": {
                  "commitsSort": ["subject", "scope"]
                }
            }
        ],
        [
            '@semantic-release/changelog',
            {
                "changelogFile": "CHANGELOG.md"
            }
        ],
        [
            '@semantic-release/git',
            {
                "message": "chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}",
                "assets": ["CHANGELOG.md"]
            }
        ],
        '@semantic-release/github'
    ]
};
