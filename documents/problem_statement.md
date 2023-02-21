## Problem statement

1) What is the problem and why it’s relevant?
2) What is your approach? (proposed solution + method)
3) How do you want to be evaluated? (deliverables + criterion for success)

## What is the problem and why it’s relevant?

Haaukins manages labs which are collections of isolated virtual machines and challenges (Docker containers) within isolated virtualnetworks.

Haaukins allocates computing and memory resources for its challenges on the local host machine. Challenges are allocated on a per-lab basis, meaning that each lab requires unique instances of given challenges. This can result in high resource usage.

At the time of writing the platform is not horizontally scalable, meaning all resources are allocated on the single host machine.

Hosting the Challenge containers on a different machine is not straightforward as there several layers of virtualization and virtual network configuration to consider.

## What is your approach?

In order to reduce resource allocation on the single host machine I want to investigate whether or not it is feasible to use modern orchestration technology and proposing changes to the platform in order to move these Challenge containers onto other hosts.

To emulate the Haaukins platform, I want to implement a minimum vertical slice that manages the labs in the same way as the platform and test to see if possible to distribute the Challenge containers and whether it is feasible.

## How do you want to be evaluated?

For deliverables, I want to submit a repository containing code and configuration for a proof-of-concept implementation of the idea proposed earlier. 

The criteria for success is based on the quality of my findings and drawn conclusions.
  